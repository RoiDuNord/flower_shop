package kafkapayment

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

func PaymentToKafka(ctx context.Context, orderQty int) {
	writer := initWriter()
	defer writer.Close()

	producer(ctx, orderQty, writer)

	fmt.Println("Payment function exiting.")
}

func initWriter() *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{"localhost:9092"},
		Topic:        "payments",
		BatchSize:    1,
		BatchTimeout: 1 * time.Millisecond,
	})
}
