package kfk

import (
	"os"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

func InitWriter(topicKey string) *kafka.Writer {
	brokersKey := "KAFKA_" + topicKey + "_BROKERS"
	topicKeyEnv := "KAFKA_" + strings.ToUpper(topicKey) + "_TOPIC"

	brokers := os.Getenv(brokersKey)
    brokerList := strings.Split(brokers, ",")

	topic := os.Getenv(topicKeyEnv)

	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:      brokerList,
		Topic:        topic,
		BatchSize:    1,
		BatchTimeout: 1 * time.Millisecond,
	})
}
