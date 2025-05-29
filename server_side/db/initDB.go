package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"server/config"

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

	database := &Database{db: db}

	if err := database.createDB(); err != nil {
		_ = database.Close()
		return nil, fmt.Errorf("failed to create tables and seed data: %w", err)
	}

	return database, nil
}

func (d *Database) Close() error {
	if err := d.db.Close(); err != nil {
		return fmt.Errorf("error closing database: %w", err)
	}
	return nil
}

func (d *Database) DB() *sql.DB {
	return d.db
}

func InitPostgres() (*Database, error) {
	dbParams, err := config.GetDBParams()
	if err != nil {
		slog.Error("error getting DB parameters", "error", err)
		return nil, err
	}

	database, err := Init(dbParams)
	if err != nil {
		slog.Error("error initializing database", "error", err)
		return nil, err
	}

	slog.Info("database initialized successfully")
	return database, nil
}
