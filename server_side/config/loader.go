package config

import (
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port     int    `yaml:"port"`
	LogLevel string `yaml:"logLevel"`
}

func loadFromFile(fileName string) (Config, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		slog.Error("error reading config file", "file", fileName, "error", err)
		return Config{}, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		slog.Error("error unmarshaling config", "file", fileName, "error", err)
		return Config{}, err
	}

	return cfg, nil
}

