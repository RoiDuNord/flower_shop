package config

import "fmt"

func (cfg *Config) validate() error {
	if cfg.Port <= 0 || cfg.Port > 65535 {
		return fmt.Errorf("неправильное значение порта (%d): допустимый диапазон 1-65535", cfg.Port)
	}
	return nil
}
