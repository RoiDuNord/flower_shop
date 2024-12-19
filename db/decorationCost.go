package db

import (
	"database/sql"
	"fmt"
	"models"

	_ "github.com/lib/pq"
)

func UpdateDecorationCost(decor models.Decoration) (int, int, error) {
	var postcardPrice int
	var err error
	postcardPriceQuery := `SELECT price FROM postcards WHERE message = $1`
	if err = db.QueryRow(postcardPriceQuery, decor.Postcard.Message).Scan(&postcardPrice); err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, fmt.Errorf("Открытка '%s' не найдена в базе данных", decor.Postcard.Message)
		}
		return 0, 0, fmt.Errorf("не удалось получить цену открытки '%s': %w", decor.Postcard.Message, err)
	}

	fmt.Printf("postcardPrice: %d\n", postcardPrice)

	var packPrice int
	flowerPriceQuery := `SELECT price FROM packs WHERE material = $1`
	if err = db.QueryRow(flowerPriceQuery, decor.Pack.Material).Scan(&packPrice); err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, fmt.Errorf("'%s' не найден в базе данных", decor.Pack.Material)
		}
		return 0, 0, fmt.Errorf("не удалось получить стоимость '%s': %w", decor.Pack.Material, err)
	}

	fmt.Printf("packPrice: %d\n", packPrice)

	return postcardPrice, packPrice, nil
}
