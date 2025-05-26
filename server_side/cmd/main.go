package main

import (
	"log/slog"
	_ "net/http/pprof"
	"server/config"
	"server/logger"
	"server/server"
)

func main() {
	logFile, err := logger.Init()
	if err != nil {
		return
	}
	defer logger.Close(logFile)
	defer slog.Info("application has been shut down")

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
