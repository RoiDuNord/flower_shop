package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"server/models"
)

type contextKey string

const PaymentDataKey contextKey = "paymentDataKey"

func (s *Server) PaymentHandling(w http.ResponseWriter, r *http.Request) {
	slog.Info("getting payment info")

	paymentData, err := readPaymentFile("payment1.json")
	if err != nil {
		slog.Error("error reading payment file", "file", "payment1.json", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	slog.Info("payment data unmarshaled successfully", "paymentData", paymentData)

	w.Header().Set("Content-Type", "application/json")

	ctx := context.WithValue(r.Context(), PaymentDataKey, paymentData)
	s.OrderHandling(w, r.WithContext(ctx))
}

func readPaymentFile(filename string) (models.Payment, error) {
	var payment models.Payment

	payInfo, err := os.ReadFile(filename)
	if err != nil {
		return payment, err
	}

	if err := json.Unmarshal(payInfo, &payment); err != nil {
		return payment, err
	}

	return payment, nil
}
