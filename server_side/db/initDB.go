package db

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func Init(params []string) (*Database, error) {
	if len(params) < 6 {
		return nil, fmt.Errorf("not enough parameters provided")
	}

	host, port, user, password, dbname, sslmode := params[0], params[1], params[2], params[3], params[4], params[5]

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	slog.Info("connected to the database")

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	if err := d.db.Close(); err != nil {
		return fmt.Errorf("error closing database: %w", err)
	}
	return nil
}
