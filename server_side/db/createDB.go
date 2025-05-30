package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"server/config"
	"server/models"
)

func createDatabaseIfNotExists(params config.DBParams) error {
	// Подключаемся к базе по умолчанию (postgres), чтобы создать новую базу
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=%s",
		params.Host, params.Port, params.User, params.Password, params.SSLMode,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	// Проверяем, есть ли база с нужным именем
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname=$1)", params.Name).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		slog.Info("database does not exist, creating", "name", params.Name)
		_, err = db.Exec("CREATE DATABASE " + params.Name)
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
	} else {
		slog.Info("database exists", "name", params.Name)
	}
	return nil
}


func (d *Database) createDB() error {
	if err := d.createTables(); err != nil {
		return err
	}

	flowers, postcards, packs := initializeData()

	if err := d.insertFlowers(flowers); err != nil {
		return err
	}
	if err := d.insertPostcards(postcards); err != nil {
		return err
	}
	if err := d.insertPacks(packs); err != nil {
		return err
	}

	slog.Info("database schema created and initial data inserted successfully")
	return nil
}

func (d *Database) createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS flowers (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100),
			color VARCHAR(50),
			cost INT,
			quantity INT
		);`,
		`CREATE TABLE IF NOT EXISTS postcards (
			id SERIAL PRIMARY KEY,
			message TEXT,
			cost INT
		);`,
		`CREATE TABLE IF NOT EXISTS packs (
			id SERIAL PRIMARY KEY,
			material VARCHAR(100),
			cost INT
		);`,
	}

	for _, query := range queries {
		if _, err := d.db.Exec(query); err != nil {
			return fmt.Errorf("error executing query: %w", err)
		}
	}
	return nil
}

func (d *Database) insertFlowers(flowers []models.Flower) error {
	for _, flower := range flowers {
		exists, err := d.existsFlower(flower.Name, flower.Color)
		if err != nil {
			return err
		}
		if exists {
			slog.Info("flower already exists, skipping insert", "name", flower.Name, "color", flower.Color)
			continue
		}

		if _, err := d.db.Exec(
			`INSERT INTO flowers (name, color, cost, quantity) VALUES ($1, $2, $3, $4)`,
			flower.Name, flower.Color, flower.Cost, flower.Quantity,
		); err != nil {
			return fmt.Errorf("inserting flower failed: %w", err)
		}
	}
	return nil
}

func (d *Database) existsFlower(name, color string) (bool, error) {
	var exists bool
	err := d.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM flowers WHERE name = $1 AND color = $2)`,
		name, color,
	).Scan(&exists)
	return exists, err
}

func (d *Database) insertPostcards(postcards []models.Postcard) error {
	for _, postcard := range postcards {
		exists, err := d.existsPostcard(postcard.Message)
		if err != nil {
			return err
		}
		if exists {
			slog.Info("postcard already exists, skipping insert", "message", postcard.Message)
			continue
		}

		if _, err := d.db.Exec(
			`INSERT INTO postcards (message, cost) VALUES ($1, $2)`,
			postcard.Message, postcard.Cost,
		); err != nil {
			return fmt.Errorf("inserting postcard failed: %w", err)
		}
	}
	return nil
}

func (d *Database) existsPostcard(message string) (bool, error) {
	var exists bool
	err := d.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM postcards WHERE message = $1)`,
		message,
	).Scan(&exists)
	return exists, err
}

func (d *Database) insertPacks(packs []models.Pack) error {
	for _, pack := range packs {
		exists, err := d.existsPack(pack.Material)
		if err != nil {
			return err
		}
		if exists {
			slog.Info("pack already exists, skipping insert", "material", pack.Material)
			continue
		}

		if _, err := d.db.Exec(
			`INSERT INTO packs (material, cost) VALUES ($1, $2)`,
			pack.Material, pack.Cost,
		); err != nil {
			return fmt.Errorf("inserting pack failed: %w", err)
		}
	}
	return nil
}

func (d *Database) existsPack(material string) (bool, error) {
	var exists bool
	err := d.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM packs WHERE material = $1)`,
		material,
	).Scan(&exists)
	return exists, err
}

func initializeData() ([]models.Flower, []models.Postcard, []models.Pack) {
	return []models.Flower{
			{Name: "Роза", Color: "Красная", Cost: 80, Quantity: 200},
			{Name: "Роза", Color: "Белая", Cost: 60, Quantity: 200},
			{Name: "Роза", Color: "Желтая", Cost: 40, Quantity: 200},
			{Name: "Лилия", Color: "Белая", Cost: 100, Quantity: 50},
			{Name: "Лилия", Color: "Желтая", Cost: 90, Quantity: 50},
			{Name: "Пион", Color: "Розовый", Cost: 120, Quantity: 100},
			{Name: "Пион", Color: "Белый", Cost: 110, Quantity: 100},
			{Name: "Лотос", Color: "Белый", Cost: 200, Quantity: 50},
			{Name: "Ромашка", Color: "Белая", Cost: 20, Quantity: 500},
		}, []models.Postcard{
			{Message: "С Днем рождения!", Cost: 5},
			{Message: "С Новым Годом!", Cost: 1},
			{Message: "Со свадьбой!", Cost: 2},
			{Message: "С Юбилеем!", Cost: 3},
			{Message: "С 8 марта!", Cost: 15},
			{Message: "С Днем Влюбленных!", Cost: 20},
		}, []models.Pack{
			{Material: "Крафт", Cost: 100},
			{Material: "Пленка", Cost: 50},
			{Material: "Лента", Cost: 10},
		}
}

// надо, чтобы создавал при входе такую бд с этими данными, если она не создана
// # PostgreSQL configuration
// DB_HOST=db
// DB_PORT=5432
// DB_USER=florist
// DB_PASSWORD=Magician1337
// DB_NAME=flower_shop
// DB_SSLMODE=disable

// package db

// import (
// 	"database/sql"
// 	"fmt"
// 	"log/slog"
// 	"server/config"

// 	_ "github.com/lib/pq"
// )

// type Database struct {
// 	db *sql.DB
// }

// func Init(params config.DBParams) (*Database, error) {
// 	psqlInfo := fmt.Sprintf(
// 		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
// 		params.Host, params.Port, params.User, params.Password, params.Name, params.SSLMode,
// 	)

// 	db, err := sql.Open("postgres", psqlInfo)
// 	if err != nil {
// 		return nil, fmt.Errorf("error opening database: %w", err)
// 	}

// 	if err = db.Ping(); err != nil {
// 		return nil, fmt.Errorf("error connecting to the database: %w", err)
// 	}

// 	slog.Info("connected to the database")

// 	database := &Database{db: db}

// 	if err := database.createDB(); err != nil {
// 		_ = database.Close()
// 		return nil, fmt.Errorf("failed to create tables and seed data: %w", err)
// 	}

// 	return database, nil
// }

// func (d *Database) Close() error {
// 	if err := d.db.Close(); err != nil {
// 		return fmt.Errorf("error closing database: %w", err)
// 	}
// 	return nil
// }

// func (d *Database) DB() *sql.DB {
// 	return d.db
// }

// func InitPostgres() (*Database, error) {
// 	dbParams, err := config.GetDBParams()
// 	if err != nil {
// 		slog.Error("error getting DB parameters", "error", err)
// 		return nil, err
// 	}

// 	database, err := Init(dbParams)
// 	if err != nil {
// 		slog.Error("error initializing database", "error", err)
// 		return nil, err
// 	}

// 	slog.Info("database initialized successfully")
// 	return database, nil
// }
