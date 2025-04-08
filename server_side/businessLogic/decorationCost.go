package businesslogic

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

	postcardCost, err := om.DB.GetDecorElementCost(PostcardTable, PostcardColumn, decor.Postcard.Message)
	if err != nil {
		slog.Error("error getting postcard cost", "error", err)
		return 0, 0, err
	}

	packCost, err := om.DB.GetDecorElementCost(PackTable, PackColumn, decor.Pack.Material)
	if err != nil {
		slog.Error("error getting pack cost", "error", err)
		return 0, 0, err
	}

	slog.Info("decoration cost retrieved", "postcardCost", postcardCost, "packCost", packCost)
	return postcardCost, packCost, nil
}
