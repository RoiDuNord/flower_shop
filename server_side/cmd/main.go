package main

import (
	"db"
	"fmt"
	"log"
	"server/order"
)

func main() {
	if err := db.InitDB(); err != nil {
		log.Fatal(err)
	}

	orderData, err := order.LoadFromFile("order1.json")
	if err != nil {
		log.Println(err)
		return
	}

	order, err := order.ParseOrder(orderData)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(string(order))

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
