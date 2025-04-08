package main

import (
	"log/slog"
	_ "net/http/pprof"
	"server/config"
	"server/logger"
	"server/server"
)

func main() {
	file := logger.Init()
	slog.Info("logger initialized")
	defer file.Close()
	defer slog.Info("shutting down application")

	cfg, err := config.ParseConfig()
	if err != nil {
		return
	}

	if cfg == (config.Config{}) {
		slog.Warn("empty config")
	}

	if err := server.Run(cfg); err != nil {
		return
	}
}
