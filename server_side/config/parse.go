package config

import (
	"fmt"
	"log/slog"
	"os"
)

func ParseConfig() (Config, error) {
	slog.Info("checking for config file", "file", "config.yaml")

	if _, err := os.Stat("config.yaml"); os.IsNotExist(err) {
		slog.Error("config file not found", "file", "config.yaml")
		return Config{}, fmt.Errorf("config.yaml not found")
	}

	cfg, err := loadFromFile("config.yaml")
	if err != nil {
		slog.Error("error loading config file", "file", "config.yaml", "error", err)
		return Config{}, err
	}

	if err := cfg.validate(); err != nil {
		slog.Error("config validation error", "error", err)
		return Config{}, err
	}

	slog.Info("config loaded and validated", "config", cfg)
	return cfg, nil
}
