package order

import (
	"encoding/json"
	"models"
)

func ParseOrder(rawOrder []byte) ([]byte, error) {
	checkedOrder, err := handleBouquetsRequest(rawOrder)
	if err != nil {
		return nil, err
	}

	orderModel, err := decodeOrder(checkedOrder)
	if err != nil {
		return nil, err
	}

	orderJSON, err := json.MarshalIndent(orderModel, "", "   ")
	if err != nil {
		return nil, err
	}

	return orderJSON, nil
}

var currentOrderID int

func decodeOrder(data []byte) (models.Order, error) {
	currentOrderID = newID()

	bouquets, err := parseBouquets(data)
	if err != nil {
		return models.Order{}, err
	}

	totalCost := calculateTotalCost(bouquets)

	order := models.Order{
		ID:        currentOrderID,
		List:      bouquets,
		OrderCost: totalCost,
	}
	return order, nil
}

func newID() int {
	currentOrderID++
	return currentOrderID
}

func calculateTotalCost(bouquets []models.Bouquet) (totalCost int) {
	for _, bouquet := range bouquets {
		totalCost += bouquet.Cost
	}
	return
}
