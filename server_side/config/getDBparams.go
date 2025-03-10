package config

import (
	"fmt"
	"log/slog"

	"github.com/joho/godotenv"
)

func GetDBParams() ([]string, error) {
	myEnv, err := godotenv.Read()
	if err != nil {
		slog.Error("error reading environment variables", "error", err)
		return nil, fmt.Errorf("error reading environment variables: %w", err)
	}

	requiredKeys := []string{"HOST", "PORT", "USER", "PASSWORD", "NAME", "SSLMODE"}
	params := make([]string, 0, len(requiredKeys))

	for _, key := range requiredKeys {
		if value, exists := myEnv[key]; !exists || value == "" {
			slog.Warn("environment variable missing or empty", "key", key)
			return nil, fmt.Errorf("missing or empty value for environment variable: %s", key)
		} else {
			params = append(params, value)
		}
	}

	return params, nil
}
