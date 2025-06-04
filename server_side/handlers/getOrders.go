package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"server/models"
	"sync"
)

func (s *Server) GetOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	slog.Info("GetOrders handler starts")
	defer slog.Info("GetOrders completed")

	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 1; i <= 5; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			order, err := fromRedis(s, i)
			if err != nil {
				slog.Error(err.Error())
				return
			}

			mu.Lock()
			defer mu.Unlock()

			if err := json.NewEncoder(w).Encode(order); err != nil {
				slog.Error("error encoding order to response", "error", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
	}

	wg.Wait()
}

func fromRedis(s *Server, numb int) (models.Order, error) {
	stringOrder, err := s.RDB.Get(s.Ctx, fmt.Sprintf("order_%d", numb)).Result()
	if err != nil {
		return models.Order{}, err
	}

	var order models.Order
	if err = json.Unmarshal([]byte(stringOrder), &order); err != nil {
		return models.Order{}, err
	}

	return order, nil
}
