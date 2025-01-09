package db

import (
	"database/sql"
	"fmt"
	"log"
)

func (d *Database) GetDecorElPrice(tableName, columnName, value string) (int, error) {
	var price int
	query := fmt.Sprintf(`SELECT price FROM %s WHERE %s = $1`, tableName, columnName)
	if err := d.db.QueryRow(query, value).Scan(&price); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("%s '%s' не найден в базе данных", columnName, value)
		}
		return 0, fmt.Errorf("не удалось получить цену '%s': %w", value, err)
	}

	log.Printf("%sPrice: %d\n", tableName[:len(tableName)-1], price)
	return price, nil
}
