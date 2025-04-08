package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"server/models"
)

func (s *Server) GetOrders(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json")

	// Log current order count before locking
	slog.Info("Current order count before locking", "count", len(s.Orders))

	s.mu.Lock()         // Lock access to Orders
	defer s.mu.Unlock() // Ensure we unlock after we're done

	// Create a slice to hold all orders
	var ordersSlice []models.Order
	for _, order := range s.Orders {
		ordersSlice = append(ordersSlice, order)
		fmt.Println(order.ID)
	}

	if err := json.NewEncoder(w).Encode(ordersSlice); err != nil {
		slog.Error("error encoding orders to response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Println("Total processed orders:", len(ordersSlice))
}
