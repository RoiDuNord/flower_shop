package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"log/slog"

	"server/config"
	"server/handlers"

	"github.com/go-chi/chi"
)

func Run(cfg config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	router := chi.NewRouter()

	s, err := handlers.NewServer(cfg, router)
	if err != nil {
		slog.Warn("running server error")
		return err
	}

	defer s.DB.Close()

	router.Get("/admin/info", s.OrderHandling)
	router.Post("/orders/payment", s.PaymentHandling)

	server := s.HTTPServer

	go func() {
		slog.Info(fmt.Sprintf("starting HTTP server on port: %d", cfg.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("error starting HTTP server", "error", err)
			return
		}
	}()

	<-ctx.Done()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("error shutting down server", "error", err)
		return err
	}

	slog.Info("server gracefully shutdown")

	return nil
}
