package order

import (
	"log"
	"os"
)

func LoadFromFile(file string) ([]byte, error) {
	orderData, err := os.ReadFile(file)
	if err != nil {
		log.Println("Ошибка при чтении информации о заказе:", err)
		return nil, err
	}
	return orderData, nil
}
