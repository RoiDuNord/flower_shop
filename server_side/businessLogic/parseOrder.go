package order

import (
	"fmt"
	"log/slog"
	"server/models"
)

func (om *OrderManager) ParseOrder(rawOrder []byte) (models.Order, error) {

	checkedOrder, err := om.handleBouquetsRequest(rawOrder)
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
	fmt.Print("currentOrderID: ")
	fmt.Println(currentOrderID)

	bouquets, err := parseBouquets(data)
	if err != nil {
		slog.Error("error parsing bouquets", "error", err)
		return models.Order{}, err
	}

	totalCost := calculateTotalCost(bouquets)
	slog.Info("total cost calculated", "totalCost", totalCost)

	// с этим кусочком надо поработать
	// paymentInfo := getPayInfo(currentOrderID)

	// isPaid := paymentInfo.IsPaid
	// if isPaid {
	// 	slog.Info("get payment", "orderID", currentOrderID)
	// }
	// с этим кусочком надо поработать

	order := models.Order{
		ID:           currentOrderID,
		BouquetsList: bouquets,
		OrderCost:    totalCost,
		// Payment: models.Payment{
		// 	IsPaid:    paymentInfo.IsPaid,
		// 	PaymentID: paymentInfo.PaymentID,
		// },
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
