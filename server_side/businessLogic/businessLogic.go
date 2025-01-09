package businesslogic

import (
	"fmt"
	"log"
	"server/config"
	"server/order"
)

func bl() {
	// params, err := config.GetDBParams()
	// if err != nil {
	// 	log.Println("ошибка получения параметров БД: %w", err)
	// }

	// db, err := db.Init(params)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }
	// defer db.Close()

	// bl - это суть работы сервера
	orderData, err := config.LoadFromFile("order1.json")
	if err != nil {
		log.Println(err)
		return
	}

	om := order.OrderManager{Db: db}
	order, err := om.ParseOrder(orderData)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(string(order))
}
