package c_order

import (
	"encoding/json"
)

var id int

func MakeOrder() ([]byte, error) {

	flowers, postcards, packs := initializeData()

	flowers1 := []Flower{flowers["redRose"], flowers["whiteLily"]}
	decoration1 := Decoration{Postcard: postcards["birthday"], Pack: packs["craft"]}
	price1 := BouquetPrice(flowers1, decoration1)

	flowers2 := []Flower{flowers["lotus"], flowers["pinkPion"]}
	decoration2 := Decoration{Postcard: postcards["happyAnniversary"], Pack: packs["film"]}
	price2 := BouquetPrice(flowers2, decoration2)

	id++
	bouquet1 := Bouquet{
		Position:   id,
		Bouquet:    flowers1,
		Price:      price1,
		Decoration: decoration1,
	}

	id++
	bouquet2 := Bouquet{
		Position:   id,
		Bouquet:    flowers2,
		Price:      price2,
		Decoration: decoration2,
	}

	bouquets := make([]Bouquet, 0)
	bouquets = append(bouquets, bouquet1, bouquet2)

	data := Bouquets{
		List: bouquets,
	}

	jsonData, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}
