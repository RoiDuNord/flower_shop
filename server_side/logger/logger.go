package logger

import (
	"log/slog"
	"os"
	"path/filepath"
)

func Init() *os.File {
	logDir, logFile := "logger", "logger.log"
	logPath := filepath.Join(logDir, logFile)

	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		slog.Error("creating log directory", "error", err)
	}

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		slog.Error("opening log file", "error", err)
	}

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}

	logger := slog.New(slog.NewJSONHandler(file, opts))
	slog.SetDefault(logger)

	return file
}
