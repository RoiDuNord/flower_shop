package order

import (
	"db"
	"encoding/json"
	"fmt"
	"log"
	"models"
)

func handleBouquetsRequest(orderData []byte) ([]byte, error) {
	bouquets, err := parseBouquets(orderData)
	if err != nil {
		return nil, err
	}

	checkedBouquets, err := updateBouquets(bouquets)
	if err != nil {
		return nil, err
	}

	for _, cal := range checkedBouquets {
		fmt.Println(cal)
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

func updateBouquets(bouquets []models.Bouquet) ([]models.Bouquet, error) {
	for i := range bouquets {
		bouquet := &bouquets[i]
		var totalCost int

		for j := range bouquet.Flowers {
			flower := &bouquet.Flowers[j]
			handledFlower, err := validateQuantityAndCost(flower)
			if err != nil {
				log.Println(err)
			}

			flower.Cost = handledFlower.Cost
			flower.Quantity = handledFlower.Quantity
			totalCost += flower.Cost
			fmt.Printf("Обновлено: %s - Количество: %d, Стоимость: %d, Общая стоимость: %d\n", flower.Name, flower.Quantity, flower.Cost, totalCost)
		}

		decorationCost, err := decorationCost(bouquet)
		if err != nil {
			log.Println(err)
			continue
		}
		totalCost += decorationCost
		bouquet.Cost = totalCost
	}
	return bouquets, nil
}

func decorationCost(bouquet *models.Bouquet) (int, error) {
	postcardPrice, packPrice, err := db.UpdateDecorationCost(bouquet.Decoration)
	if err != nil {
		return 0, err
	}
	bouquet.Decoration.Postcard.Price = postcardPrice
	bouquet.Decoration.Pack.Price = packPrice
	bouquet.Decoration.Cost = postcardPrice + packPrice
	return bouquet.Decoration.Cost, nil
}

func validateQuantityAndCost(flower *models.Flower) (*models.Flower, error) {
	fullFlowerName := fmt.Sprintf("%s %s", flower.Name, flower.Color)

	availQty, cost, err := db.UpdateQuantityAndCost(flower.Name, flower.Color, flower.Quantity)
	log.Printf("Количество для %s обновлено на %d", fullFlowerName, availQty)

	flower.Cost, flower.Quantity = cost, availQty

	return flower, err
}
