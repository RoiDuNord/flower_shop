package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"log/slog"

	order "server/businessLogic"
	"server/config"
	"server/db"
	"server/models"

	"github.com/go-chi/chi"
)

type Server struct {
	DB         *db.Database
	HTTPServer *http.Server
	Context    context.Context
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

func NewServer(cfg config.Config, router *chi.Mux) (*Server, error) {
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
	}

	slog.Info("server created successfully", "address", s.HTTPServer.Addr)

	return s, nil
}

func (s *Server) OrderHandling(w http.ResponseWriter, r *http.Request) {
	slog.Info("received request for InfoHandler")

	orderData, err := order.ProcessOrder(s.DB)
	if err != nil {
		slog.Error("error processing order", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if len(orderData.BouquetsList) == 0 {
		slog.Warn("no order found")
		http.Error(w, "No order found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if payment, ok := r.Context().Value(PaymentDataKey).(models.Payment); ok {
		slog.Info("payment data key found", "payment", payment)
		orderData.PaymentID = payment.PaymentID
	}

	if err := json.NewEncoder(w).Encode(orderData); err != nil {
		slog.Error("encoding orderData problem", "error", err)
		http.Error(w, "Failed to encode order data", http.StatusInternalServerError)
		return
	}

	slog.Info("order successfully returned", "order", orderData)
}
