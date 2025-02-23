package order

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"server/models"
)

func (om *OrderManager) ParseOrder(rawOrder []byte) ([]byte, error) {
	checkedOrder, err := om.handleBouquetsRequest(rawOrder)
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

	paymentInfo := getPayInfo(currentOrderID)

	isPaid := paymentInfo.IsPaid

	order := models.Order{
		ID:        currentOrderID,
		List:      bouquets,
		OrderCost: totalCost,
		Payment:   isPaid,
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

func getPayInfo(ID int) models.Payment {
	fmt.Println(ID)

	payInfo, err := os.ReadFile("payment1.json")
	if err != nil {
		log.Println(err)
		return models.Payment{}
	}

	fmt.Println(payInfo)

	var payment models.Payment

	if err := json.Unmarshal(payInfo, &payment); err != nil {
		fmt.Println(err)
		return models.Payment{}
	}
	fmt.Println(payment)

	return payment
}
