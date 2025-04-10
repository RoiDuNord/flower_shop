package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"log/slog"

	bl "server/businessLogic"
	"server/config"
	"server/db"
	"server/models"
	ko "server/services/kafkaOrder"
	kp "server/services/kafkaPayment"

	"github.com/go-chi/chi"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	DB         *db.Database
	RDB        *redis.Client
	HTTPServer *http.Server
	Context    context.Context
}

func initPostgres() (*db.Database, error) {
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

func initRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "1234",
	})
}

func NewServer(cfg config.Config, router *chi.Mux, ctx context.Context) (*Server, error) {
	database, err := initPostgres()
	if err != nil {
		return nil, err
	}

	rDB := initRedis()

	s := &Server{
		DB:  database,
		RDB: rDB,
		HTTPServer: &http.Server{
			Addr:    fmt.Sprintf("localhost:%d", cfg.Port),
			Handler: router,
		},
		Context: ctx,
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
	payments := make(map[int]models.Payment)

	var wg sync.WaitGroup

	wg.Add(4)
	go s.consumeOrders(&wg, orderChan)
	go s.consumePayments(&wg, paymentChan)
	go processPayments(&wg, paymentChan, payments)
	go s.combineOrders(&wg, orderChan, payments, w)

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
	return make(chan models.Order, orderQty), make(chan models.Payment, orderQty)
}

func (s *Server) consumeOrders(wg *sync.WaitGroup, orderChan chan models.Order) {
	defer wg.Done()
	slog.Info("consumeOrders started")
	defer slog.Info("consumeOrders finished")

	start := time.Now()

	ko.Consumer(s.Context, orderChan, cap(orderChan))

	duration := time.Since(start)
	fmt.Printf("consumeOrders took %s\n", duration)
}

func (s *Server) consumePayments(wg *sync.WaitGroup, paymentChan chan models.Payment) {
	defer wg.Done()

	slog.Info("consumePayments started")
	defer slog.Info("consumePayments finished")

	start := time.Now()

	kp.Consumer(s.Context, paymentChan, cap(paymentChan))

	duration := time.Since(start)
	fmt.Printf("consumePayments took %s\n", duration)
}

func processPayments(wg *sync.WaitGroup, paymentChan chan models.Payment, payments map[int]models.Payment) {
	defer wg.Done()
	for payment := range paymentChan {
		payments[payment.OrderID] = payment
	}
}

func (s *Server) combineOrders(wg *sync.WaitGroup, orderChan chan models.Order, payments map[int]models.Payment, w http.ResponseWriter) {
	defer wg.Done()
	for rawOrder := range orderChan {
		s.processOrder(rawOrder, payments, w)
	}
}

func (s *Server) processOrder(rawOrder models.Order, payments map[int]models.Payment, w http.ResponseWriter) {
	processedOrder, err := bl.ProcessOrder(s.Context, s.DB, rawOrder)
	if err != nil {
		slog.Error(fmt.Sprintf("error processing order: %e", err))
		return
	}

	if err := handlePayment(processedOrder.ID, payments, &processedOrder); err != nil {
		slog.Error(err.Error())
		return
	}

	if err := s.addProcessedOrder(processedOrder); err != nil {
		slog.Warn(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	slog.Info("order successfully processed", "order", processedOrder)
}

func handlePayment(orderID int, payments map[int]models.Payment, processedOrder *models.Order) error {
	payment, exists := payments[orderID]
	if !exists {
		return fmt.Errorf("payment not found for order ID: %d", orderID)
	}
	addPayment(payment, processedOrder)
	return nil
}

func (s *Server) addProcessedOrder(processedOrder models.Order) error {
	if len(processedOrder.BouquetsList) == 0 {
		return fmt.Errorf("no bouquets found in order")
	}

	binaryOrder, err := json.Marshal(processedOrder)
	if err != nil {
		slog.Error("error marshalling processed order", "error", err)
		return err
	}

	if err := s.RDB.Set(s.Context, fmt.Sprintf("order_%d", processedOrder.ID), binaryOrder, 0).Err(); err != nil {
		slog.Error("error adding order to Redis", "error", err)
		return err
	}
	slog.Info("added order to Redis", "orderID", processedOrder.ID)

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
