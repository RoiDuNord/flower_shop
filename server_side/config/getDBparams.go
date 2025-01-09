package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadFromFile(file string) ([]byte, error) {
	orderData, err := os.ReadFile(file)
	if err != nil {
		log.Println("Ошибка при чтении информации о заказе:", err)
		return nil, err
	}
	return orderData, nil
}

func GetDBParams() ([]string, error) {
	myEnv, err := godotenv.Read()
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения переменных окружения: %w", err)
	}

	requiredKeys := []string{"HOST", "PORT", "USER", "PASSWORD", "NAME", "SSLMODE"}
	params := make([]string, 0, len(requiredKeys))

	for _, key := range requiredKeys {
		if value, exists := myEnv[key]; !exists || value == "" {
			return nil, fmt.Errorf("не установлено значение для переменной окружения: %s", key)
		} else {
			params = append(params, value)
		}
	}

	return params, nil
}
