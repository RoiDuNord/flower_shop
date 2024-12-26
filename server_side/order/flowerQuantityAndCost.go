package order

import (
	"fmt"
)

func (db *DBWrapper) UpdateQuantityAndCost(flowerName, flowerColor string, requestedQuantity int) (int, int, error) {
	flower := fmt.Sprintf("%s %s", flowerName, flowerColor)

	availableQuantity, price, err := db.GetFlowersQtyAndPrice(flowerName, flowerColor)
	if err != nil {
		return 0, 0, err
	}

	quantity, err := validateQuantity(requestedQuantity, availableQuantity, flower)

	if err := db.UpdateQuantity(quantity, flowerName, flowerColor); err != nil {
		return 0, 0, fmt.Errorf("не удалось обновить количество для '%s': %w", flower, err)
	}

	totalCost := price * quantity
	return quantity, totalCost, err
}

func validateQuantity(requestedQuantity, availableQuantity int, flower string) (int, error) {
	if requestedQuantity > availableQuantity {
		return availableQuantity, fmt.Errorf("Максимум %d для %s", availableQuantity, flower)
	}
	return requestedQuantity, nil
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
