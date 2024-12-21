package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() error {

	params, err := getDBParams()
	if err != nil {
		log.Fatalf("Ошибка получения параметров БД: %v", err)
	}
	fmt.Println(params)

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		params[0], params[1], params[2], params[3], params[4])

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return fmt.Errorf("ошибка при открытии базы данных: %w", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("ошибка при подключении к БД: %w", err)
	}
	fmt.Println("Успешное подключение к БД!")
	return nil
}

func getDBParams() ([]string, error) {
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

func Close() error {

	if err := db.Close(); err != nil {
		return fmt.Errorf("ошибка при закрытии базы данных: %w", err)
	}

	return nil
}
