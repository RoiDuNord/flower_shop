package handlers

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"log/slog"

	bl "server/businessLogic"
	"server/config"
	"server/db"
	"server/models"
	ko "server/services/kafkaOrder"
	kp "server/services/kafkaPayment"

	"github.com/go-chi/chi"
)

type Server struct {
	DB         *db.Database
	HTTPServer *http.Server
	Context    context.Context
	Orders     map[int]models.Order
	mu         sync.Mutex
}

func newDB() (*db.Database, error) {
	dbParams, err := config.GetDBParams()
	if err != nil {
		slog.Error("error getting DB parameters", "error", err)
		return nil, err
	}

	database, err := db.Init(dbParams)
	if err != nil {
		slog.Error("error initializing database", "error", err)
		return nil, err
	}

	slog.Info("database initialized successfully")
	return database, nil
}

func NewServer(cfg config.Config, router *chi.Mux, ctx context.Context) (*Server, error) {
	database, err := newDB()
	if err != nil {
		return nil, err
	}

	s := &Server{
		DB: database,
		HTTPServer: &http.Server{
			Addr:    fmt.Sprintf("localhost:%d", cfg.Port),
			Handler: router,
		},
		Context: ctx,
		Orders:  make(map[int]models.Order, 5),
	}

	slog.Info("server created successfully", "address", s.HTTPServer.Addr)

	return s, nil
}

func (s *Server) CreateOrder(w http.ResponseWriter, r *http.Request) {
	slog.Info("received request for CreateOrder")

	orderQty, err := s.getOrderQuantity()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	orderChan, paymentChan := initializeChannels(orderQty)
	orders, payments := initializeMaps()

	var wg sync.WaitGroup

	wg.Add(4)
	go s.consumeOrders(&wg, orderChan)     // сначала это
	go s.consumePayments(&wg, paymentChan) // потом это
	go processPayments(&wg, paymentChan, payments)
	go s.combineOrders(&wg, orderChan, payments, orders, w)

	wg.Wait()
}

func (s *Server) getOrderQuantity() (int, error) {
	orderQty, ok := s.Context.Value("orderQty").(int)
	if !ok {
		slog.Error("error retrieving order quantity from context")
		return 0, fmt.Errorf("invalid order quantity")
	}
	return orderQty, nil
}

func initializeChannels(orderQty int) (chan models.Order, chan models.Payment) {
	orderChan := make(chan models.Order, orderQty)
	paymentChan := make(chan models.Payment, orderQty)
	return orderChan, paymentChan
}

func initializeMaps() (map[int]models.Order, map[int]models.Payment) {
	orders := make(map[int]models.Order)
	payments := make(map[int]models.Payment)
	return orders, payments
}

func (s *Server) consumeOrders(wg *sync.WaitGroup, orderChan chan models.Order) {
	defer wg.Done()
	ko.Consumer(s.Context, orderChan, cap(orderChan))
}

func (s *Server) consumePayments(wg *sync.WaitGroup, paymentChan chan models.Payment) {
	defer wg.Done()
	kp.Consumer(s.Context, paymentChan, cap(paymentChan))
}

func processPayments(wg *sync.WaitGroup, paymentChan chan models.Payment, payments map[int]models.Payment) {
	defer wg.Done()
	for payment := range paymentChan {
		payments[payment.OrderID] = payment
	}
}

func (s *Server) combineOrders(wg *sync.WaitGroup, orderChan chan models.Order, payments map[int]models.Payment, orders map[int]models.Order, w http.ResponseWriter) {
	defer wg.Done()
	for rawOrder := range orderChan {
		s.processOrder(rawOrder, payments, orders, w)
	}
}

func (s *Server) processOrder(rawOrder models.Order, payments map[int]models.Payment, orders map[int]models.Order, w http.ResponseWriter) {
	processedOrder, err := bl.ProcessOrder(s.Context, s.DB, rawOrder)
	if err != nil {
		slog.Error(fmt.Sprintf("error processing order: %e", err))
		return
	}

	if err := s.handlePayment(processedOrder.ID, payments, &processedOrder); err != nil {
		slog.Error(err.Error())
		return
	}

	if err := s.addProcessedOrder(processedOrder, orders); err != nil {
		slog.Warn(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	slog.Info("order successfully processed", "order", processedOrder)
}

func (s *Server) handlePayment(orderID int, payments map[int]models.Payment, processedOrder *models.Order) error {
	payment, exists := payments[orderID]
	if !exists {
		return fmt.Errorf("payment not found for order ID: %d", orderID)
	}
	addPayment(payment, processedOrder)
	return nil
}

func (s *Server) addProcessedOrder(processedOrder models.Order, orders map[int]models.Order) error {
	if len(processedOrder.BouquetsList) == 0 {
		return fmt.Errorf("no bouquets found in order")
	}

	orders[processedOrder.ID] = processedOrder
	s.Orders[processedOrder.ID] = processedOrder
	return nil
}

func addPayment(payment models.Payment, order *models.Order) {
	if payment.IsPaid {
		order.PaymentID = payment.PaymentID
		order.PaymentStatus = "Оплачено"
	} else {
		order.PaymentStatus = "Не оплачено"
	}
}
