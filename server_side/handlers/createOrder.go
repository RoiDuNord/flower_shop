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
	rds "server/redis"
	ko "server/services/kafkaOrder"
	kp "server/services/kafkaPayment"

	"github.com/go-chi/chi"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	DB         *db.Database
	RDB        *redis.Client
	HTTPServer *http.Server
	Ctx        context.Context
	Router     *chi.Mux
}

func NewServer(cfg config.Config, ctx context.Context) (*Server, error) {
	database, err := db.InitPostgres()
	if err != nil {
		return nil, err
	}

	rDB, err := rds.InitRedis()
	if err != nil {
		return nil, err
	}

	router := chi.NewRouter()

	s := &Server{
		DB:     database,
		RDB:    rDB,
		Router: router,
		Ctx:    ctx,
	}

	s.HTTPServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: router,
	}

	slog.Info("server created successfully", "address", s.HTTPServer.Addr)

	return s, nil
}

func (s *Server) CloseAllDBs() {
	if err := s.RDB.Close(); err != nil {
		slog.Error("error closing Redis client", "error", err)
	}
	defer slog.Info("redis has been successfully shut down")
	if err := s.DB.Close(); err != nil {
		slog.Error("error closing database", "error", err)
	}
	defer slog.Info("database has been successfully shut down")
}

func (s *Server) CreateOrder(w http.ResponseWriter, r *http.Request) {
	slog.Info("received request for CreateOrder")
	defer slog.Info("CreateOrder completed")

	orderQty, ok := s.Ctx.Value(models.OrderQtyKey{}).(int)
	if ok {
		slog.Info("get order quantity from context", "quantity:", orderQty)
	} else {
		slog.Error("order quantity not found in context")
	}

	orderChan, paymentChan := initChannels(orderQty)
	payments := make(map[int]models.Payment)

	var wg sync.WaitGroup

	wg.Add(4)
	go s.consumeOrders(&wg, orderChan)
	go s.consumePayments(&wg, paymentChan)
	go processPayments(&wg, paymentChan, payments)
	go s.combineOrders(&wg, orderChan, payments, w)

	wg.Wait()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"all orders successfully processed"}`))
}

func initChannels(orderQty int) (chan models.Order, chan models.Payment) {
	return make(chan models.Order, orderQty), make(chan models.Payment, orderQty)
}

func (s *Server) consumeOrders(wg *sync.WaitGroup, orderChan chan models.Order) {
	cfg, err := config.GetKafkaParams("ORDERS")
	if err != nil {
		slog.Error("failed to load kafka configs", "error", err)
	}

	defer wg.Done()
	slog.Info("consumeOrders started")
	defer slog.Info("consumeOrders finished")

	start := time.Now()

	ko.Consumer(s.Ctx, cfg, orderChan, cap(orderChan))

	duration := time.Since(start)
	fmt.Printf("consumeOrders took %s\n", duration)
}

func (s *Server) consumePayments(wg *sync.WaitGroup, paymentChan chan models.Payment) {
	cfg, err := config.GetKafkaParams("PAYMENTS")
	if err != nil {
		slog.Error("failed to load kafka configs", "error", err)
	}

	defer wg.Done()

	slog.Info("consumePayments started")
	defer slog.Info("consumePayments finished")

	start := time.Now()

	kp.Consumer(s.Ctx, cfg, paymentChan, cap(paymentChan))

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
	processedOrder, err := bl.ProcessOrder(s.Ctx, s.DB, rawOrder)
	if err != nil {
		slog.Error("error processing order", "error", err)
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

	if err := s.RDB.Set(s.Ctx, fmt.Sprintf("order_%d", processedOrder.ID), binaryOrder, 0).Err(); err != nil {
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
