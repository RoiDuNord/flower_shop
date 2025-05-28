package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log/slog"

	"server/config"
	"server/handlers"
	"server/models"
	ko "server/services/kafkaOrder"
	kp "server/services/kafkaPayment"

	"github.com/go-chi/chi"
)

func Run(cfg config.Config) error {
	ctx, cancel := initContext()
	defer cancel()

	orderQty := 5
	ctx = context.WithValue(ctx, models.OrderQtyKey{}, orderQty)

	s, err := handlers.NewServer(cfg, ctx)
	if err != nil {
		slog.Warn("running server error")
		return err
	}
	defer s.CloseAllDBs()

	setupRoutes(s, s.Router)

	go startHTTPServer(s.HTTPServer, cfg.Port)

	sendOrdersAndPaymentsToKafka(ctx, orderQty)

	return shutdownServer(s)
}

func initContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 400*time.Second)
}

func setupRoutes(s *handlers.Server, router *chi.Mux) {
	router.Post("/orders/create", s.CreateOrder)
	router.Get("/orders/get", s.GetOrders)
}

func startHTTPServer(server *http.Server, port int) {
	slog.Info(fmt.Sprintf("starting HTTP server on port: %d", port))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("error starting HTTP server", "error", err)
	}
}

func sendOrdersAndPaymentsToKafka(ctx context.Context, orderQty int) {
	ko.SendOrderToKafka(ctx, orderQty)
	kp.SendPaymentToKafka(ctx, orderQty)
}

func shutdownServer(s *handlers.Server) error {
	server := s.HTTPServer
	shutdownSignals := make(chan os.Signal, 1)
	signal.Notify(shutdownSignals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	select {
	case <-shutdownSignals:
		slog.Info("received shutdown signal")
	case <-s.Ctx.Done():
		slog.Info("context deadline exceeded")
	}

	if err := server.Shutdown(s.Ctx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)

		if err := server.Close(); err != nil {
			slog.Error("forced shutdown failed", "error", err)
			return err
		}
	}

	slog.Info("server shutdown complete")
	return nil
}