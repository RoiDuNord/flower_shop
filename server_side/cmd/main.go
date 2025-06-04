package main

import (
	"log/slog"
	_ "net/http/pprof"
	"server/config"
	"server/logger"
	"server/server"
)

func main() {
	cfg, err := config.ParseConfig()
	if err != nil {
		return
	}

	if cfg == (config.Config{}) {
		slog.Warn("empty config")
	}

	logFile, err := logger.Init(cfg.LogLevel)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer logger.Close(logFile)
	defer slog.Info("application has been shut down")

	if err := server.Run(cfg); err != nil {
		return
	}

}
