package config

import (
	"fmt"
	"log/slog"

	"github.com/joho/godotenv"
)


func getParams(prefix string, requiredKeys []string) ([]string, error) {
	myEnv, err := readEnv()
	if err != nil {
		return nil, err
	}

	envKeys := buildKeys(prefix, requiredKeys)
	params, err := checkEnvValues(myEnv, envKeys)
	if err != nil {
		return nil, err
	}

	return params, nil
}


func readEnv() (map[string]string, error) {
	myEnv, err := godotenv.Read()
	if err != nil {
		slog.Error("error reading environment variables", "error", err)
		return nil, fmt.Errorf("error reading environment variables: %w", err)
	}
	return myEnv, nil
}

func buildKeys(prefix string, requiredKeys []string) []string {
	keys := make([]string, len(requiredKeys))

	for i, key := range requiredKeys {
		keys[i] = prefix + "_" + key
	}

	return keys
}

func checkEnvValues(myEnv map[string]string, envKeys []string) ([]string, error) {
	params := make([]string, 0, len(envKeys))

	for _, key := range envKeys {
		if value, exists := myEnv[key]; !exists || value == "" {
			slog.Warn("environment variable missing or empty", "key", key)
			return nil, fmt.Errorf("missing or empty value for environment variable: %s", key)
		} else {
			params = append(params, value)
		}
	}

	return params, nil
}