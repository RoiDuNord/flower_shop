package order

import (
	"fmt"
	"server/models"
)

func (om *OrderManager) GetFlowerAvailQtyAndCost(flower *models.Flower) error {
	availQty, price, err := om.Db.GetFlowerQtyAndPrice(flower.Name, flower.Color)
	if err != nil {
		return err
	}

	validateQuantity(flower, availQty)

	if err := om.Db.UpdateFlowerQty(flower.Quantity, flower.Name, flower.Color); err != nil {
		return fmt.Errorf("не удалось обновить количество для '%s': %w", flower.Name+" "+flower.Color, err)
	}

	flower.Cost = price * flower.Quantity
	return err
}

func validateQuantity(flower *models.Flower, availableQuantity int) {
	if flower.Quantity > availableQuantity {
		flower.Quantity = availableQuantity
		flower.ErrorMessage = fmt.Sprintf("доступно %d шт.", availableQuantity)
	}
}
