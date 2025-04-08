package businesslogic

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"server/models"
)

func (om *OrderManager) handleBouquetsRequest(ctx context.Context, rawOrderData models.Order) ([]byte, error) {
	checkedBouquets, err := om.updateBouquets(rawOrderData.BouquetsList)
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

func (om *OrderManager) updateBouquets(bouquets []models.Bouquet) ([]models.Bouquet, error) {
	for i := range bouquets {
		bouquet := &bouquets[i]

		for j := range bouquet.FlowerList {
			flower := &bouquet.FlowerList[j]
			if err := om.validateFlowersQtyAndCost(flower, &bouquet.BouquetCost); err != nil {
				slog.Error("error validating flower", "flower", flower.Name+" "+flower.Color, "error", err)
			}
			fullFlowerName := fmt.Sprintf("%s %s", flower.Name, flower.Color)
			slog.Info("flower updated", "name", fullFlowerName, "quantity", flower.Quantity, "cost", flower.Cost)
			updateBouquetCost(om, bouquet)
		}

		// if bouquet.BouquetCost != 0 {
		// 	if err := om.decorationCost(bouquet, &bouquet.BouquetCost); err != nil {
		// 		slog.Error("error updating decoration cost for bouquet", "error", err)
		// 	}
		// } else {
		// 	slog.Error("zero-price bouquet")
		// }

		fmt.Printf("%d bouquet total cost: %d\n", i+1, bouquet.BouquetCost)
	}

	return bouquets, nil
}

func updateBouquetCost(om *OrderManager, bouquet *models.Bouquet) {
	if bouquet.BouquetCost == 0 {
		slog.Error("zero-price bouquet")
		return
	}

	if err := om.decorationCost(bouquet, &bouquet.BouquetCost); err != nil {
		slog.Error("error updating decoration cost for bouquet", "error", err)
	}
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

	slog.Info(fmt.Sprintf("'%s' - quantity: %d, cost: %d", fullFlowerName, flower.Quantity, flower.Cost))

	updateTotalCost(totalCost, flower.Cost)

	return nil
}

func updateTotalCost(totalCost *int, amount int) {
	*totalCost += amount
	slog.Debug("total cost updated", "newTotalCost", *totalCost)
}
