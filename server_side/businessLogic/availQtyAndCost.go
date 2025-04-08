package businesslogic

import (
	"fmt"
	"log/slog"
	"server/models"
)

func (om *OrderManager) GetFlowerAvailQtyAndCost(flower *models.Flower) error {
	flowerName := fmt.Sprintf("%s %s", flower.Name, flower.Color)

	slog.Info(fmt.Sprintf("getting %s availability and cost", flowerName))

	availQty, cost, err := om.DB.GetFlowerQtyAndCost(flower.Name, flower.Color)
	if err != nil {
		slog.Error("error getting flower quantity and cost", "error", err)
		return err
	}

	slog.Info("available quantity and cost retrieved", "availableQuantity", availQty, "cost", cost)

	validateQuantity(flower, availQty)

	if err := om.DB.UpdateFlowerQty(flower.Quantity, flower.Name, flower.Color); err != nil {
		slog.Error("error updating flower quantity", "flower", flowerName, "error", err)
		return fmt.Errorf("failed to update quantity for '%s': %w", flowerName, err)
	}

	flower.Cost = cost * flower.Quantity
	// slog.Info("flower cost calculated", "flowerName", flowerName, "cost", flower.Cost)
	return nil
}

func validateQuantity(flower *models.Flower, availableQuantity int) {
	flowerName := fmt.Sprintf("%s %s", flower.Name, flower.Color)
	slog.Info("validating quantity", "flower", flowerName, "orderedQuantity", flower.Quantity, "availableQuantity", availableQuantity)

	if flower.Quantity > availableQuantity {
		flower.Quantity = availableQuantity
		flower.ErrorMessage = fmt.Sprintf("доступно %d шт.", availableQuantity)
		slog.Warn("quantity adjusted", "flower", flowerName, "adjustedQuantity", flower.Quantity)
	}
}
