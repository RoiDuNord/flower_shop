package order

import (
	"fmt"
	"log"
	"server/models"
)

func (om *OrderManager) GetFlowerAvailQtyAndCost(flower *models.Flower) error {
	availQty, price, err := om.Db.GetFlowerQtyAndPrice(flower.Name, flower.Color)
	if err != nil {
		return err
	}

	validateQuantity(flower, availQty)

	if err := om.Db.UpdateQty(flower.Quantity, flower.Name, flower.Color); err != nil {
		return fmt.Errorf("не удалось обновить количество для '%s': %w", flower.Name+" "+flower.Color, err)
	}

	flower.Cost = price * flower.Quantity
	log.Println(flower.Cost)
	return err
}

func validateQuantity(flower *models.Flower, availableQuantity int) {
	if flower.Quantity > availableQuantity {
		flower.Quantity = availableQuantity
		flower.ErrorMessage = fmt.Sprintf("доступно %d шт.", availableQuantity)
	}
}

// if requestedQuantity > availableQuantity {
// 	return availableQuantity, availableQuantity * fmt.Errorf("Максимум %d\n для '%s'", availableQuantity, flower)
// }

// requestedQuantity = availableQuantity

// var updatedQuantity int
// query3 := `SELECT quantity FROM flowers WHERE name = $1 AND color = $2`
// if err = db.QueryRow(query3, flowerName, flowerColor).Scan(&updatedQuantity); err != nil {
// 	if err == sql.ErrNoRows {
// 		return 0, fmt.Errorf("'%s' нет в базе данных", flower)
// 	}
// 	return 0, fmt.Errorf("не удалось получить количество '%s: %w", flower, err)
// }

// fmt.Printf("updatedQuantity: %d\n", updatedQuantity)
