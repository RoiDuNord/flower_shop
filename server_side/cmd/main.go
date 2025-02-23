package main

import (
	"log"
	"server/config"
	"server/server"
)

func main() {
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}

	if cfg == (config.Config{}) {
		log.Fatal("Empty config")
	}

	server.Run(cfg)
}

// func newLogger() (*slog.Logger, error) {
// 	logDir, logFile := "logger", "sysLog.log"
// 	logPath := filepath.Join(logDir, logFile)

// 	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
// 	if err != nil {
// 		slog.Error("opening log file", "error", err)
// 		return nil, err
// 	}
// 	defer file.Close()

// 	opts := &slog.HandlerOptions{
// 		AddSource: true,
// 		Level:     slog.LevelDebug,
// 	}

// 	logger := slog.New(slog.NewJSONHandler(file, opts))

// 	slog.SetDefault(logger)

// 	return logger, nil
// }
