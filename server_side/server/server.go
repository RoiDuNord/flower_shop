package server

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"server/config"
	"server/db"
	"server/models"

	"github.com/go-chi/chi"
)

type Server struct {
	db *db.Database
}

func Run(cfg config.Config) error {
	// этот кусок
	logDir, logFile := "logger", "sysLog.log"
	logPath := filepath.Join(logDir, logFile)

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		slog.Error("opening log file", "error", err)
		return err
	}
	defer file.Close()

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}

	Logger := slog.New(slog.NewJSONHandler(file, opts))

	slog.SetDefault(Logger)
	// до этого момента

	router := chi.NewRouter()

	s, err := newServer()
	if err != nil {
		log.Println(err)
		return err
	}

	router.Get("/admin/info", s.InfoHandler)
	router.Post("/payOrder", s.PayOrder)

	slog.Info("starting HTTP server on port")
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), router); err != nil {
		slog.Info("starting HTTP server error")
	}

	// data := bl.OrderManage(s.Db)

	// fmt.Println(data)

	return nil
}

func newDB() (*db.Database, error) {
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

func newServer() (*Server, error) {

	database, err := newDB()
	if err != nil {
		slog.String("error", err.Error())
		return nil, err
	}

	s := &Server{
		db: database,
	}

	return s, nil
}

func (s *Server) InfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Server info"))
}

func (s *Server) PayOrder(w http.ResponseWriter, r *http.Request) {
	var payStatus models.Payment
	if err := json.NewDecoder(r.Body).Decode(&payStatus); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	defer r.Body.Close()

	order := models.Order{}
	log.Printf("new factor %v", order.Payment)
}
