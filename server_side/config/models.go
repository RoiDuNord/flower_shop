package config

type DBParams struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisParams struct {
	Host     string
	Port     string
	Password string
}

type KafkaParams struct {
	Brokers string
	Topic   string
	GroupID string
}
