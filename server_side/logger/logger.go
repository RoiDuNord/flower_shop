package logger

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type logger = *os.File

func Init(logLevel string) (logger, error) {
	logDir, logFile := "logger", "logger.log"
	logPath := filepath.Join(logDir, logFile)

	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}

	msk, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return nil, err
	}

	opts := &slog.HandlerOptions{
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				formattedTime := a.Value.Time().In(msk).Format("2006-01-02 15:04:05")
				return slog.String(slog.TimeKey, formattedTime)
			}
			return a
		},
		Level: parseLogLevel(logLevel),
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

func parseLogLevel(level string) slog.Level {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
