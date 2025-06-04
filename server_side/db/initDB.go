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

func Init(params config.DBParams) (*Database, error) {
	if err := createDatabaseIfNotExists(params); err != nil {
		return nil, err
	}

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		params.Host, params.Port, params.User, params.Password, params.Name, params.SSLMode,
	)
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
