package kafkaorder

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

func SendOrderToKafka(ctx context.Context, orderQty int) {
	writer := initWriter()
	defer writer.Close()

	producer(ctx, orderQty, writer)

	fmt.Println("Order function exiting.")
}

func initWriter() *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{"localhost:9092"},
		Topic:        "orders",
		BatchSize:    1,
		BatchTimeout: 1 * time.Millisecond,
	})
}

// go consumer(ctx, orderCh, orderQty)

// printCh(orderCh)

// func printCh(orderCh chan models.Order) {
// 	for order := range orderCh {
// 		fmt.Println("Received order in main:", order)
// 	}
// }
