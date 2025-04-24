package config

func GetDBParams() ([]string, error) {
	return getParams("POSTGRES", []string{"HOST", "PORT", "USER", "PASSWORD", "NAME", "SSLMODE"})
}

func GetRedisParams() ([]string, error) {
	return getParams("REDIS", []string{"HOST", "PORT", "PASSWORD"})
}
