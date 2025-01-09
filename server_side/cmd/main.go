package main

import (
	"log"
	"server/config"
	"server/server"
)

func main() {
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}

	if cfg == (config.Config{}) {
		log.Fatal("Empty config")
	}

	server.Run(cfg)
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
