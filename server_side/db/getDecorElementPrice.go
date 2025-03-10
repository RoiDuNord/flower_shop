package db

import (
	"database/sql"
	"fmt"
	"log/slog"
)

func (d *Database) GetDecorElementPrice(tableName, columnName, value string) (int, error) {
	var price int
	query := fmt.Sprintf(`SELECT price FROM %s WHERE %s = $1`, tableName, columnName)

	if err := d.db.QueryRow(query, value).Scan(&price); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("%s '%s' not found in the database", columnName, value)
		}
		return 0, fmt.Errorf("could not get price '%s': %w", value, err)
	}

	slog.Info(fmt.Sprintf("%s price: %d", tableName[:len(tableName)-1], price))
	return price, nil
}
