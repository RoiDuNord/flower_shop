package config

import "fmt"

func (cfg *Config) validate() error {
	if cfg.Port <= 0 || cfg.Port > 65535 {
		return fmt.Errorf("invalid port value (%d): acceptable range [1:65535]", cfg.Port)
	}
	return nil
}
