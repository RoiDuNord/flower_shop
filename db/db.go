package db

import (
	"database/sql"
	"fmt"
)

var db *sql.DB

const (
	Host     = "localhost"
	Port     = 5432
	User     = "florist"
	Password = "Magician1337"
	Name     = "flower_shop"
)

func InitDB(host string, port int, user, password, dbname string) error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
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

func Close() error {
	if db != nil {
		if err := db.Close(); err != nil {
			return fmt.Errorf("ошибка при закрытии базы данных: %w", err)
		}
	}
	return nil
}
