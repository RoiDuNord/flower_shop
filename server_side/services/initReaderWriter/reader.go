package kfk

import (
	"os"
	"strings"

	"github.com/segmentio/kafka-go"
)

func InitReader(topicKey string) *kafka.Reader {
	brokersKey := "KAFKA_" + strings.ToUpper(topicKey) + "_BROKERS"
	topicKeyEnv := "KAFKA_" + strings.ToUpper(topicKey) + "_TOPIC"
	groupIDKey := "KAFKA_" + strings.ToUpper(topicKey) + "_GROUP_ID"

	brokers := os.Getenv(brokersKey)
	if brokers == "" {
		brokers = "localhost:9092"
	}
	brokerList := strings.Split(brokers, ",")

	topic := os.Getenv(topicKeyEnv)
	if topic == "" {
		topic = strings.ToLower(topicKey)
	}

	groupID := os.Getenv(groupIDKey)
	if groupID == "" {
		groupID = strings.ToLower(topicKey) + "-consumer-group"
	}

	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokerList,
		Topic:   topic,
		GroupID: groupID,
	})
}
