package order

import (
	"server/models"
)

const (
	PostcardTable  = "postcards"
	PostcardColumn = "message"
	PackTable      = "packs"
	PackColumn     = "material"
)

func (db *DBWrapper) GetDecorationCost(decor models.Decoration) (int, int, error) {
	postcardPrice, err := db.GetDecorElPrice(PostcardTable, PostcardColumn, decor.Postcard.Message)
	if err != nil {
		return 0, 0, err
	}

	packPrice, err := db.GetDecorElPrice(PackTable, PackColumn, decor.Pack.Material)
	if err != nil {
		return 0, 0, err
	}

	return postcardPrice, packPrice, nil
}
