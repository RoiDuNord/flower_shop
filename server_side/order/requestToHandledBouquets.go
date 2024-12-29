package order

import (
	"encoding/json"
	"fmt"
	"log"

	"server/db"
	"server/models"
)

type OrderManager struct {
	Db *db.Database
}

func (om *OrderManager) handleBouquetsRequest(orderData []byte) ([]byte, error) {
	bouquets, err := parseBouquets(orderData)
	if err != nil {
		return nil, err
	}

	checkedBouquets, err := om.updateBouquets(bouquets)
	if err != nil {
		return nil, err
	}

	handledOrder, err := json.MarshalIndent(checkedBouquets, "", "   ")
	if err != nil {
		return nil, err
	}

	return handledOrder, nil
}

func parseBouquets(orderData []byte) ([]models.Bouquet, error) {
	var bouquets []models.Bouquet
	if err := json.Unmarshal(orderData, &bouquets); err != nil {
		return nil, err
	}
	return bouquets, nil
}

func (om *OrderManager) updateBouquets(bouquets []models.Bouquet) ([]models.Bouquet, error) {
	for i := range bouquets {
		bouquet := &bouquets[i]

		for j := range bouquet.Flowers {
			flower := &bouquet.Flowers[j]
			if err := om.validateFlowersQtyAndCost(flower, &bouquet.Cost); err != nil {
				log.Printf("Ошибка при валидации цветка %s: %v", flower.Name, err)
			}
			log.Println(flower)
			fmt.Printf("Обновлено: %s - Количество: %d, Стоимость: %d, Общая стоимость: %d\n", flower.Name, flower.Quantity, flower.Cost, bouquet.Cost)
		}

		if err := om.decorationCost(bouquet, &bouquet.Cost); err != nil {
			log.Printf("Ошибка при обновлении стоимости дополнений для букета: %v", err)
		}
	}

	return bouquets, nil
}

func (om *OrderManager) decorationCost(bouquet *models.Bouquet, totalCost *int) error {
	postcardPrice, packPrice, err := om.GetDecorationCost(bouquet.Decoration)
	if err != nil {
		return err
	}
	bouquet.Decoration.Postcard.Price = postcardPrice
	bouquet.Decoration.Pack.Price = packPrice
	bouquet.Decoration.Cost = postcardPrice + packPrice

	updateTotalCost(totalCost, bouquet.Decoration.Cost)

	return nil
}

func (om *OrderManager) validateFlowersQtyAndCost(flower *models.Flower, totalCost *int) error {
	if err := om.GetFlowerAvailQtyAndCost(flower); err != nil {
		return err
	}

	fullFlowerName := fmt.Sprintf("%s %s", flower.Name, flower.Color)
	log.Printf("Количество для %s обновлено на %d", fullFlowerName, flower.Quantity)

	log.Printf("Qty: %d, C: %d", flower.Quantity, flower.Cost)

	updateTotalCost(totalCost, flower.Cost)

	return nil
}

func updateTotalCost(totalCost *int, amount int) {
	*totalCost += amount
}

// func (db *OrderManager) validateQuantityAndCost(flower *models.Flower) (*models.Flower, error) {
// 	fullFlowerName := fmt.Sprintf("%s %s", flower.Name, flower.Color)

// 	availQty, cost, err := db.GetFlowerAvailQtyAndCost(flower.Name, flower.Color, flower.Quantity)
// 	log.Printf("Количество для %s обновлено на %d", fullFlowerName, availQty)

// 	flower.Cost, flower.Quantity = cost, availQty

// 	return flower, err
// }
