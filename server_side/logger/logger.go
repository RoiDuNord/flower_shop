package logger

import (
	"log/slog"
	"os"
	"path/filepath"
)

type logger = *os.File

func Init() (logger, error) {
	logDir, logFile := "logger", "logger.log"
	logPath := filepath.Join(logDir, logFile)

	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}

	opts := &slog.HandlerOptions{
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				formattedTime := a.Value.Time().Format("15:04:05")
				return slog.String(slog.TimeKey, formattedTime)
			}
			return a
		},
	}
	logger := slog.New(slog.NewJSONHandler(file, opts))
	slog.SetDefault(logger)

	logger.Info("logger initialized successfully", slog.String("module", "logger"))
	return file, nil
}

func Close(logger logger) error {
	if logger != nil {
		if err := logger.Close(); err != nil {
			slog.Error("failed to close logger", "error", err)
			return err
		}
		slog.Info("logger closed successfully")
	}
	return nil
}
