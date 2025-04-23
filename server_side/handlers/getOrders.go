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

	fmt.Println("order.ID", order.ID)

	return order, nil
}

// package handlers

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log/slog"
// 	"net/http"
// 	"server/models"
// 	"sync"
// )

// func (s *Server) GetOrders(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	slog.Info("GetOrders handler starts")

// 	var (
// 		wg sync.WaitGroup
// 		mu sync.Mutex
// 	)

// 	wg.Add(5)
// 	for i := 1; i <= 5; i++ {
// 		go func(i int) {
// 			defer wg.Done()
// 			mu.Lock()
// 			if err := fromRedis(w, s, &mu, i); err != nil {
// 				slog.Error(err.Error())
// 				return
// 			}
// 			mu.Unlock()
// 		}(i)
// 	}

// 	// fmt.Println("Total processed orders:", len(ordersSlice))
// 	wg.Wait()
// }

// func fromRedis(w http.ResponseWriter, s *Server, mu *sync.Mutex, numb int) error {
// 	binaryOrder, err := s.RDB.Get(s.Context, fmt.Sprintf("user_%d", numb)).Result()
// 	if err != nil {
// 		return err
// 	}

// 	var order models.Order
// 	err = json.Unmarshal([]byte(binaryOrder), &order)
// 	if err != nil {
// 		return err
// 	}

// 	if err := json.NewEncoder(w).Encode(order); err != nil {
// 		// slog.Error("error encoding orders to response", "error", err)
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		return err
// 	}

// 	fmt.Println("order.ID", order.ID)

// 	return nil
// }
