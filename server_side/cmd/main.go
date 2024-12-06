package main

import (
	c_order "client/makeOrder"
	"encoding/json"
	"fmt"
)

func main() {
	orderInfo, err := c_order.MakeOrder()
	if err != nil {
		fmt.Println("Ошибка при передаче данных:", err)
		return
	}

	order, err := Parse(orderInfo)
	if err != nil {
		fmt.Println("Ошибка при парсинге:", err)
		return
	}

	orderJSON, _ := json.MarshalIndent(order, "", "   ")
	fmt.Println(string(orderJSON))
}

var id int

type Order struct {
	ID         int               `json:"orderID"`
	List       []c_order.Bouquet `json:"bouquetsList"`
	TotalPrice int               `json:"orderPrice"`
}

func Parse(data []byte) (Order, error) {
	id = newID()

	var bouquets c_order.Bouquets
	if err := json.Unmarshal(data, &bouquets); err != nil {
		return Order{}, err
	}

	var totalPrice int
	for _, bouquet := range bouquets.List {
		totalPrice += bouquet.Price
	}

	order := Order{
		ID:         id,
		List:       bouquets.List,
		TotalPrice: totalPrice,
	}
	return order, nil
}

func newID() int {
	id++
	return id
}
