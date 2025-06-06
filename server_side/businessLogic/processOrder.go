package businesslogic

import (
	"context"
	"encoding/json"
	"log/slog"
	"server/db"
	"server/models"
)

type OrderManager struct {
	DB *db.Database
}

func ProcessOrder(ctx context.Context, db *db.Database, rawOrderData models.Order) (models.Order, error) {
	om := OrderManager{DB: db}
	order, err := om.ParseOrder(ctx, rawOrderData)
	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (om *OrderManager) ParseOrder(ctx context.Context, rawOrderData models.Order) (models.Order, error) {
	checkedOrder, err := om.handleBouquetsRequest(rawOrderData)
	if err != nil {
		slog.Error("error handling bouquets request", "error", err)
		return models.Order{}, err
	}

	order, err := decodeOrder(checkedOrder)
	if err != nil {
		slog.Error("error decoding order", "error", err)
		return models.Order{}, err
	}

	slog.Info("order parsed successfully", "orderID", order.ID)
	return order, nil
}

var currentOrderID int

func decodeOrder(data []byte) (models.Order, error) {
	currentOrderID = newID()

	bouquets, err := parseBouquets(data)
	if err != nil {
		slog.Error("error parsing bouquets", "error", err)
		return models.Order{}, err
	}

	totalCost := calculateTotalCost(bouquets)
	slog.Info("total cost calculated", "totalCost", totalCost)

	order := models.Order{
		ID:           currentOrderID,
		BouquetsList: bouquets,
		OrderCost:    totalCost,
	}

	return order, nil
}

func newID() int {
	currentOrderID++
	return currentOrderID
}

func calculateTotalCost(bouquets []models.Bouquet) (totalCost int) {
	for _, bouquet := range bouquets {
		totalCost += bouquet.BouquetCost
	}
	return
}

func parseBouquets(orderData []byte) ([]models.Bouquet, error) {
	var bouquets []models.Bouquet
	if err := json.Unmarshal(orderData, &bouquets); err != nil {
		return nil, err
	}
	return bouquets, nil
}
