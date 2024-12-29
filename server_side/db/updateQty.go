package db

import (
	"fmt"
)

func (d *Database) UpdateQty(quantity int, flowerName, flowerColor string) error {
	query := `UPDATE flowers 
		SET quantity = quantity - $1 
		WHERE name = $2 AND color = $3`
	if _, err := d.db.Exec(query, quantity, flowerName, flowerColor); err != nil {
		return fmt.Errorf("не удалось выполнить запрос обновления: %w", err)
	}
	return nil
}
