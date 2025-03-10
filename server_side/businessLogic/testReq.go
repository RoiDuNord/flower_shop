package order

import (
	"log/slog"
	"os"
	"server/db"
	"server/models"
)

func loadFromFile(file string) ([]byte, error) {
	orderData, err := os.ReadFile(file)
	if err != nil {
		slog.Error("error reading order data file", "error", err)
		return nil, err
	}
	return orderData, nil
}

func ProcessOrder(db *db.Database) (models.Order, error) {
	orderData, err := loadFromFile("order1.json")
	if err != nil {
		return models.Order{}, err
	}

	om := OrderManager{DB: db}
	order, err := om.ParseOrder(orderData)
	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}
