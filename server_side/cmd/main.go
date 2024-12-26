package main

import (
	"fmt"
	"log"
	"server/config"
	"server/db"
	"server/order"
)

func main() {
	params, err := config.GetDBParams()
	if err != nil {
		fmt.Errorf("ошибка получения параметров БД: %w", err)
	}

	db, err := db.NewDB(params)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// bl - это суть работы сервера
	orderData, err := order.LoadFromFile("order1.json")
	if err != nil {
		log.Println(err)
		return
	}

	order, err := db.ParseOrder(orderData)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(string(order))
	//
}

// logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
// if err != nil {
// 	fmt.Println("Ошибка при открытии файла логов:", err)
// 	return
// }
// defer logFile.Close()
// log.SetOutput(logFile)

// orderInfo, err := c_order.MakeOrder()
// if err != nil {
// 	log.Println("Ошибка при передаче данных:", err)
// 	return
// }
// fmt.Println(string(orderInfo))
