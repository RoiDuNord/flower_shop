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
	} // дальше этой строки не идет

	slog.Info("here") //

	if cfg == (config.Config{}) {
		slog.Warn("empty config")
	}

	slog.Info("here 2")

	logFile, err := logger.Init(cfg.LogLevel)
	if err != nil {
		return
	}
	defer logger.Close(logFile)
	defer slog.Info("application has been shut down")

	slog.Info("here 3") // только сюда не доходит

	slog.Info("before server.Run")
	if err := server.Run(cfg); err != nil {
		slog.Error("server.Run returned error", "error", err)
	}
	slog.Info("after server.Run")
}
