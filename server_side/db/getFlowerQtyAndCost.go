package db

import (
	"database/sql"
	"fmt"
	"log/slog"
)

func (d *Database) GetFlowerQtyAndCost(flowerName, flowerColor string) (int, int, error) {
	var availableQuantity, cost int
	getQtyAndCostQuery := `SELECT quantity, cost
		FROM flowers
		WHERE name = $1 AND color = $2`

	if err := d.db.QueryRow(getQtyAndCostQuery, flowerName, flowerColor).Scan(&availableQuantity, &cost); err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, fmt.Errorf("no data found for '%s %s'", flowerName, flowerColor)
		}
		return 0, 0, fmt.Errorf("could not retrieve data for '%s %s': %w", flowerName, flowerColor, err)
	}

	slog.Info(fmt.Sprintf("retrieved %d units at cost %d for %s %s", availableQuantity, cost, flowerName, flowerColor))
	return availableQuantity, cost, nil
}
