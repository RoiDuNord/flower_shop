package config

func Parse() (Config, error) {
	cfg, err := loadFromFile("config.yaml")
	if err != nil {
		return Config{}, err
	}

	if err := cfg.validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
