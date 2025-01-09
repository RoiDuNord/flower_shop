package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

// type UpdaterDB interface {
// 	UpdateQuantityAndCost() error
// 	UpdateDecorationCost() error
// } куда можно?

func Init(params []string) (*Database, error) {
	host, port, user, password, dbname, sslmode := params[0], params[1], params[2], params[3], params[4], params[5]

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("ошибка при открытии базы данных: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка при подключении к БД: %w", err)
	}

	fmt.Println("Успешное подключение к БД!")

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	if err := d.db.Close(); err != nil {
		return fmt.Errorf("ошибка при закрытии базы данных: %w", err)
	}
	return nil
}
