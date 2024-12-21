package db

import (
	"database/sql"
	"fmt"
	"log"
	"models"
)

const (
	PostcardTable  = "postcards"
	PostcardColumn = "message"
	PackTable      = "packs"
	PackColumn     = "material"
)

func UpdateDecorationCost(decor models.Decoration) (int, int, error) {
	postcardPrice, err := getPrice(PostcardTable, PostcardColumn, decor.Postcard.Message)
	if err != nil {
		return 0, 0, err
	}

	packPrice, err := getPrice(PackTable, PackColumn, decor.Pack.Material)
	if err != nil {
		return 0, 0, err
	}

	return postcardPrice, packPrice, nil
}

func getPrice(tableName, columnName, value string) (int, error) {
	if value == "" {
		return 0, fmt.Errorf("%s не должно быть пустым", columnName)
	}

	var price int
	query := fmt.Sprintf(`SELECT price FROM %s WHERE %s = $1`, tableName, columnName)
	if err := db.QueryRow(query, value).Scan(&price); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("%s '%s' не найден в базе данных", columnName, value)
		}
		return 0, fmt.Errorf("не удалось получить цену '%s': %w", value, err)
	}

	log.Printf("%sPrice: %d\n", columnName, price)
	return price, nil
}
