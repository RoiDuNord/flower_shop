package db

import (
	"fmt"
)

func UpdateQuantityAndCost(flowerName, flowerColor string, requestedQuantity int) (int, int, error) {
	flower := fmt.Sprintf("%s %s", flowerName, flowerColor)

	availableQuantity, price, err := getQtyAndPrice(flowerName, flowerColor)
	if err != nil {
		return 0, 0, err
	}

	quantity, err := validateQuantity(requestedQuantity, availableQuantity, flower)

	if err := updateQuantity(quantity, flowerName, flowerColor); err != nil {
		return 0, 0, fmt.Errorf("не удалось обновить количество для '%s': %w", flower, err)
	}

	totalCost := price * quantity
	return quantity, totalCost, err
}

func getQtyAndPrice(flowerName, flowerColor string) (int, int, error) {
	if db == nil {
		return 0, 0, fmt.Errorf("база данных не инициализирована")
	}

	var availableQuantity, price int
	getQtyAndPriceQuery := `SELECT quantity, price
		FROM flowers
		WHERE name = \$1 AND color = \$2`
	if err := db.QueryRow(getQtyAndPriceQuery, flowerName, flowerColor).Scan(&availableQuantity, &price); err != nil {
		return 0, 0, fmt.Errorf("не удалось получить данные для '%s %s': %w", flowerName, flowerColor, err)
	}
	return availableQuantity, price, nil
}

func validateQuantity(requestedQuantity, availableQuantity int, flower string) (int, error) {
	if requestedQuantity > availableQuantity {
		return availableQuantity, fmt.Errorf("Максимум %d для %s", availableQuantity, flower)
	}
	return requestedQuantity, nil
}

func updateQuantity(quantity int, flowerName, flowerColor string) error {
	query := `
		UPDATE flowers 
		SET quantity = quantity - $1 
		WHERE name = $2 AND color = $3
	`
	if _, err := db.Exec(query, quantity, flowerName, flowerColor); err != nil {
		return fmt.Errorf("не удалось выполнить запрос обновления: %w", err)
	}
	return nil
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
