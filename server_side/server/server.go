// package server

// import (
// 	"context"
// 	"fmt"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"syscall"
// 	"time"

// 	"log/slog"

// 	"server/config"
// 	"server/handlers"
// 	"server/models"
// 	ko "server/services/kafkaOrder"
// 	kp "server/services/kafkaPayment"

// 	"github.com/go-chi/chi"
// )

// // Run initializes the server and starts handling requests
// func Run(cfg config.Config) error {
// 	ctx, cancel := initContext()
// 	defer cancel()

// 	orderQty := 5
// 	ctx = context.WithValue(ctx, models.OrderQtyKey{}, orderQty)

// 	s, err := handlers.NewServer(cfg, ctx)
// 	if err != nil {
// 		slog.Warn("running server error")
// 		return err
// 	}
// 	defer s.Close()

// 	setupRoutes(s, s.Router)

// 	go startHTTPServer(s.HTTPServer, cfg.Port)

// 	sendOrdersAndPaymentsToKafka(ctx, orderQty)

// 	return shutdownServer(s)
// }

// func initContext() (context.Context, context.CancelFunc) {
// 	return context.WithTimeout(context.Background(), 5*time.Second)
// }

// func setupRoutes(s *handlers.Server, router *chi.Mux) {
// 	router.Post("/orders/create", s.CreateOrder)
// 	router.Get("/orders/get", s.GetOrders)
// }

// func startHTTPServer(server *http.Server, port int) {
// 	slog.Info(fmt.Sprintf("starting HTTP server on port: %d", port))
// 	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 		slog.Error("error starting HTTP server", "error", err)
// 	}
// }

// func sendOrdersAndPaymentsToKafka(ctx context.Context, orderQty int) {
// 	ko.SendOrderToKafka(ctx, orderQty)
// 	kp.SendPaymentToKafka(ctx, orderQty)
// }

// func shutdownServer(s *handlers.Server) error {
// 	server := s.HTTPServer
// 	shutdown := make(chan os.Signal, 1)
// 	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

// 	<-shutdown
// 	slog.Info("Start shutdown...")

// 	// Attempt graceful shutdown
// 	if err := server.Shutdown(s.Ctx); err != nil {
// 		slog.Error("graceful shutdown failed", "error", err)

// 		// Attempt forced shutdown if graceful shutdown fails
// 		if err := server.Close(); err != nil {
// 			slog.Error("forced shutdown failed", "error", err)
// 			return err
// 		}
// 	}

// 	slog.Info("Server shutdown complete")
// 	return nil
// }



package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"log/slog"

	"server/config"
	"server/handlers"
	"server/models"
	ko "server/services/kafkaOrder"
	kp "server/services/kafkaPayment"
)

func Run(cfg config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	orderQty := 5

	ctx = context.WithValue(ctx, models.OrderQtyKey{}, orderQty)

	s, err := handlers.NewServer(cfg, ctx)
	if err != nil {
		slog.Warn("running server error")
		return err
	}

	defer s.DB.Close()

	server := s.HTTPServer

	s.Router.Post("/orders/create", s.CreateOrder)
	s.Router.Get("/orders/get", s.GetOrders)

	go func() {
		slog.Info(fmt.Sprintf("starting HTTP server on port: %d", cfg.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("error starting HTTP server", "error", err)
			return
		}
	}()

	ko.SendOrderToKafka(ctx, orderQty)
	kp.SendPaymentToKafka(ctx, orderQty)

	<-ctx.Done()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("error shutting down server", "error", err)
		return err
	}

	slog.Info("server gracefully shutdown")

	return nil
}
