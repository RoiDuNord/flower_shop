package order

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"server/db"
	"server/models"
)

type OrderManager struct {
	DB *db.Database
}

func (om *OrderManager) handleBouquetsRequest(orderData []byte) ([]byte, error) {
	bouquets, err := parseBouquets(orderData)
	if err != nil {
		slog.Error("error parsing bouquets", "error", err)
		return nil, err
	}

	checkedBouquets, err := om.updateBouquets(bouquets)
	if err != nil {
		slog.Error("error updating bouquets", "error", err)
		return nil, err
	}

	handledOrder, err := json.MarshalIndent(checkedBouquets, "", "   ")
	if err != nil {
		slog.Error("error marshaling checked bouquets to json", "error", err)
		return nil, err
	}

	slog.Info("bouquets request handled successfully")
	return handledOrder, nil
}

func parseBouquets(orderData []byte) ([]models.Bouquet, error) {
	var bouquets []models.Bouquet
	if err := json.Unmarshal(orderData, &bouquets); err != nil {
		return nil, err
	}
	return bouquets, nil
}

func (om *OrderManager) updateBouquets(bouquets []models.Bouquet) ([]models.Bouquet, error) {
	for i := range bouquets {
		bouquet := &bouquets[i]

		for j := range bouquet.FlowerList {
			flower := &bouquet.FlowerList[j]
			if err := om.validateFlowersQtyAndCost(flower, &bouquet.BouquetCost); err != nil {
				slog.Error("error validating flower", "name", flower.Name, "error", err)
			}
			slog.Info("flower updated", "name", flower.Name, "quantity", flower.Quantity, "cost", flower.Cost)
			fmt.Printf("updated: %s - quantity: %d, cost: %d\n", flower.Name, flower.Quantity, flower.Cost)
		}
		fmt.Printf("%d bouquet total cost: %d\n", i+1, bouquet.BouquetCost)

		if err := om.decorationCost(bouquet, &bouquet.BouquetCost); err != nil {
			slog.Error("error updating decoration cost for bouquet", "error", err)
		}
	}

	return bouquets, nil
}

func (om *OrderManager) decorationCost(bouquet *models.Bouquet, totalCost *int) error {
	postcardPrice, packPrice, err := om.GetDecorationCost(bouquet.Decoration)
	if err != nil {
		return err
	}
	bouquet.Decoration.Postcard.Cost = postcardPrice
	bouquet.Decoration.Pack.Cost = packPrice
	bouquet.Decoration.DecorationCost = postcardPrice + packPrice

	updateTotalCost(totalCost, bouquet.Decoration.DecorationCost)

	return nil
}

func (om *OrderManager) validateFlowersQtyAndCost(flower *models.Flower, totalCost *int) error {
	if err := om.GetFlowerAvailQtyAndCost(flower); err != nil {
		return err
	}

	fullFlowerName := fmt.Sprintf("%s %s", flower.Name, flower.Color)
	slog.Info("quantity updated", "flower", fullFlowerName, "quantity", flower.Quantity)

	slog.Info("flower cost info", "quantity", flower.Quantity, "cost", flower.Cost)

	updateTotalCost(totalCost, flower.Cost)

	return nil
}

func updateTotalCost(totalCost *int, amount int) {
	*totalCost += amount
	slog.Debug("total cost updated", "newTotalCost", *totalCost)
}
