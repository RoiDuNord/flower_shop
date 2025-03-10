package order

import (
	"log/slog"
	"server/models"
)

const (
	PostcardTable  = "postcards"
	PostcardColumn = "message"
	PackTable      = "packs"
	PackColumn     = "material"
)

func (om *OrderManager) GetDecorationCost(decor models.Decoration) (int, int, error) {
	slog.Info("getting decoration cost", "postcard", decor.Postcard.Message, "pack", decor.Pack.Material)

	postcardPrice, err := om.DB.GetDecorElementPrice(PostcardTable, PostcardColumn, decor.Postcard.Message)
	if err != nil {
		slog.Error("error getting postcard price", "error", err)
		return 0, 0, err
	}

	packPrice, err := om.DB.GetDecorElementPrice(PackTable, PackColumn, decor.Pack.Material)
	if err != nil {
		slog.Error("error getting pack price", "error", err)
		return 0, 0, err
	}

	slog.Info("decoration cost retrieved", "postcardPrice", postcardPrice, "packPrice", packPrice)
	return postcardPrice, packPrice, nil
}
