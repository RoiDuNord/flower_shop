package db

import (
	"fmt"
)

func (d *Database) GetFlowerQtyAndPrice(flowerName, flowerColor string) (int, int, error) {
	var availableQuantity, price int
	getQtyAndPriceQuery := `SELECT quantity, price
		FROM flowers
		WHERE name = $1 AND color = $2`
	if err := d.db.QueryRow(getQtyAndPriceQuery, flowerName, flowerColor).Scan(&availableQuantity, &price); err != nil {
		return 0, 0, fmt.Errorf("не удалось получить данные для '%s %s': %w", flowerName, flowerColor, err)
	}
	return availableQuantity, price, nil
}
