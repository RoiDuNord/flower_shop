package config

import "strings"

func GetDBParams() ([]string, error) {
	return getParams("DB", []string{"HOST", "PORT", "USER", "PASSWORD", "NAME", "SSLMODE"})
}

func GetRedisParams() ([]string, error) {
	return getParams("REDIS", []string{"HOST", "PORT", "PASSWORD"})
}

func GetKafkaParams(name string) ([]string, error) {
	return getParams("KAFKA_"+strings.ToUpper(name), []string{"BROKERS", "TOPIC", "GROUP_ID"})
}
