package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"log/slog"

	"server/config"
	"server/handlers"
	ko "server/services/kafkaOrder"
	kp "server/services/kafkaPayment"

	"github.com/go-chi/chi"
)

func Run(cfg config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	orderQty := 5
	ctx = context.WithValue(ctx, "orderQty", orderQty)

	router := chi.NewRouter()

	s, err := handlers.NewServer(cfg, router, ctx)
	if err != nil {
		slog.Warn("running server error")
		return err
	}

	defer s.DB.Close()

	server := s.HTTPServer

	router.Post("/orders/create", s.CreateOrder)
	router.Get("/orders/get", s.GetOrders)
	// router.Post("/orders/payment/process", s.CheckAndAttachPayment)

	go func() {
		slog.Info(fmt.Sprintf("starting HTTP server on port: %d", cfg.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("error starting HTTP server", "error", err)
			return
		}
	}()

	ko.OrderToKafka(ctx, orderQty)
	kp.PaymentToKafka(ctx, orderQty)

	<-ctx.Done()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("error shutting down server", "error", err)
		return err
	}

	slog.Info("server gracefully shutdown")

	return nil
}
