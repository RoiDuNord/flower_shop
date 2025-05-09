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
	Ctx        context.Context
	Router     *chi.Mux
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

func initRedis() (*redis.Client, error) {
	rdb, err := redisClient()
	if err != nil {
		slog.Error("error initializing Redis client", "error", err)
	}

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		slog.Error("error connecting to Redis", "error", err)
		return nil, err
	}

	slog.Info("redis client initialized successfully")
	return rdb, nil
}

func redisClient() (*redis.Client, error) {
	redisParams, err := config.GetRedisParams()
	if err != nil {
		slog.Error("error getting Redis parameters", "error", err)
		return nil, err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisParams[0], redisParams[1]),
		Password: fmt.Sprintf("%s", redisParams[2]),
	})

	return rdb, nil
}

func NewServer(cfg config.Config, ctx context.Context) (*Server, error) {
	database, err := initPostgres()
	if err != nil {
		return nil, err
	}

	rDB, err := initRedis()
	if err != nil {
		return nil, err
	}

	router := chi.NewRouter()

	s := &Server{
		DB:  database,
		RDB: rDB,
		HTTPServer: &http.Server{
			Addr:    fmt.Sprintf("localhost:%d", cfg.Port),
			Handler: router,
		},
		Router: router,
		Ctx:    ctx,
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

func initializeChannels(orderQty int) (chan models.Order, chan models.Payment) {
	return make(chan models.Order, orderQty), make(chan models.Payment, orderQty)
}

func (s *Server) consumeOrders(wg *sync.WaitGroup, orderChan chan models.Order) {
	defer wg.Done()
	slog.Info("consumeOrders started")
	defer slog.Info("consumeOrders finished")

	start := time.Now()

	ko.Consumer(s.Ctx, orderChan, cap(orderChan))

	duration := time.Since(start)
	fmt.Printf("consumeOrders took %s\n", duration)
}

func (s *Server) consumePayments(wg *sync.WaitGroup, paymentChan chan models.Payment) {
	defer wg.Done()

	slog.Info("consumePayments started")
	defer slog.Info("consumePayments finished")

	start := time.Now()

	kp.Consumer(s.Ctx, paymentChan, cap(paymentChan))

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
