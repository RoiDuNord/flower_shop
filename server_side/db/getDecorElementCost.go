package db

import (
	"database/sql"
	"fmt"
)

func (d *Database) GetDecorElementCost(tableName, columnName, value string) (int, error) {
	var cost int
	query := fmt.Sprintf(`SELECT cost FROM %s WHERE %s = $1`, tableName, columnName)

	if err := d.db.QueryRow(query, value).Scan(&cost); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("%s '%s' not found in the database", columnName, value)
		}
		return 0, fmt.Errorf("could not get cost '%s': %w", value, err)
	}

	return cost, nil
}
