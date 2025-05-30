package kfk

import (
	"server/config"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

func InitReader(params config.KafkaParams) *kafka.Reader {
	brokerList := strings.Split(params.Brokers, ",")

	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokerList,
		Topic:   params.Topic,
		GroupID: params.GroupID,
	})
}

func InitWriter(params config.KafkaParams) *kafka.Writer {
	brokerList := strings.Split(params.Brokers, ",")

	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:      brokerList,
		Topic:        params.Topic,
		BatchSize:    1,
		BatchTimeout: time.Millisecond,
	})
}
