package db

import (
	"fmt"
	"log/slog"
)

func (d *Database) UpdateFlowerQty(quantity int, flowerName, flowerColor string) error {
	query := `UPDATE flowers 
		SET quantity = quantity - $1
		WHERE name = $2 AND color = $3`

	result, err := d.db.Exec(query, quantity, flowerName, flowerColor)
	if err != nil {
		return fmt.Errorf("could not execute updated query: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not retrieve affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no flower with name '%s' and color '%s'", flowerName, flowerColor)
	}

	slog.Info(fmt.Sprintf("updated quantity by %d for %s %s", quantity, flowerName, flowerColor))
	return nil
}
