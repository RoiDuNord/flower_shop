package config

import (
	"strings"
)

func GetDBParams() (DBParams, error) {
	params, err := getParams("DB", []string{"HOST", "PORT", "USER", "PASSWORD", "NAME", "SSLMODE"})
	if err != nil {
		return DBParams{}, err
	}
	return DBParams{
		Host:     params["HOST"],
		Port:     params["PORT"],
		User:     params["USER"],
		Password: params["PASSWORD"],
		Name:     params["NAME"],
		SSLMode:  params["SSLMODE"],
	}, nil
}

func GetRedisParams() (RedisParams, error) {
	params, err := getParams("REDIS", []string{"HOST", "PORT", "PASSWORD"})
	if err != nil {
		return RedisParams{}, err
	}
	return RedisParams{
		Host:     params["HOST"],
		Port:     params["PORT"],
		Password: params["PASSWORD"],
	}, nil
}

func GetKafkaParams(name string) (KafkaParams, error) {
	params, err := getParams("KAFKA_"+strings.ToUpper(name), []string{"BROKERS", "TOPIC", "GROUP_ID"})
	if err != nil {
		return KafkaParams{}, err
	}
	return KafkaParams{
		Brokers: params["BROKERS"],
		Topic:   params["TOPIC"],
		GroupID: params["GROUP_ID"],
	}, nil
}
