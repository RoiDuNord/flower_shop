package db

import (
	"database/sql"
	"fmt"
	"log/slog"
)

func (d *Database) GetFlowerQtyAndPrice(flowerName, flowerColor string) (int, int, error) {
	var availableQuantity, price int
	getQtyAndPriceQuery := `SELECT quantity, price
		FROM flowers
		WHERE name = $1 AND color = $2`

	if err := d.db.QueryRow(getQtyAndPriceQuery, flowerName, flowerColor).Scan(&availableQuantity, &price); err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, fmt.Errorf("no data found for '%s %s'", flowerName, flowerColor)
		}
		return 0, 0, fmt.Errorf("could not retrieve data for '%s %s': %w", flowerName, flowerColor, err)
	}

	slog.Info(fmt.Sprintf("retrieved %d units at price %d for %s %s", availableQuantity, price, flowerName, flowerColor))
	return availableQuantity, price, nil
}
