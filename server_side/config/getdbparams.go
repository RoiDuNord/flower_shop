package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

func GetDBParams() ([]string, error) {
	myEnv, err := godotenv.Read()
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения переменных окружения: %w", err)
	}

	requiredKeys := []string{"HOST", "PORT", "USER", "PASSWORD", "NAME"}
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
