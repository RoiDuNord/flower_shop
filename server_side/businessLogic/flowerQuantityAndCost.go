package order

import (
	"fmt"
	"log/slog"
	"server/models"
)

func (om *OrderManager) GetFlowerAvailQtyAndCost(flower *models.Flower) error {
	flowerName := fmt.Sprintf("%s %s", flower.Name, flower.Color)

	slog.Info("getting flower availability and cost", "flowerName", flowerName)

	availQty, price, err := om.DB.GetFlowerQtyAndPrice(flower.Name, flower.Color)
	if err != nil {
		slog.Error("error getting flower quantity and price", "error", err)
		return err
	}

	slog.Info("available quantity and price retrieved", "availableQuantity", availQty, "price", price)

	validateQuantity(flower, availQty)

	if err := om.DB.UpdateFlowerQty(flower.Quantity, flower.Name, flower.Color); err != nil {
		slog.Error("error updating flower quantity", "flower", flowerName, "error", err)
		return fmt.Errorf("failed to update quantity for '%s': %w", flowerName, err)
	}

	flower.Cost = price * flower.Quantity
	slog.Info("flower cost calculated", "flowerName", flowerName, "cost", flower.Cost)
	return nil
}

func validateQuantity(flower *models.Flower, availableQuantity int) {
	slog.Info("validating quantity", "flower", flower.Name+" "+flower.Color, "currentQuantity", flower.Quantity, "availableQuantity", availableQuantity)

	if flower.Quantity > availableQuantity {
		flower.Quantity = availableQuantity
		flower.ErrorMessage = fmt.Sprintf("доступно %d шт.", availableQuantity)
		slog.Warn("quantity adjusted", "flower", flower.Name+" "+flower.Color, "adjustedQuantity", flower.Quantity)
	}
}
