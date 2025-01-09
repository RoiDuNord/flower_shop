package server

import (
	"fmt"
	"log"
	"net/http"
	"server/config"
	"server/db"

	"github.com/go-chi/chi"
)

type Server struct {
	db *db.Database
}

func Run(cfg config.Config) error {
	router := chi.NewRouter()

	s, err := NewServer()
	if err != nil {
		log.Println(err)
		return err
	}

	fmt.Println(s)

	// router.Get("/admin/info", s.InfoHandler)

	log.Printf("Starting HTTP server on port %d", cfg.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), router); err != nil {
		log.Fatalf("Error starting HTTP server: %v", err)
	}

	return nil
}

func NewDB() (*db.Database, error) {
	dbParams, err := config.GetDBParams()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	database, err := db.Init(dbParams)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return database, nil
}

func NewServer() (*Server, error) {
	database, err := NewDB()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	s := &Server{
		db: database,
	}

	return s, nil
}
