package kafkaorder

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"server/models"

	"github.com/segmentio/kafka-go"
)

func Consumer(ctx context.Context, orderChan chan models.Order, unprocessedQty int) {
	reader := initReader()
	defer reader.Close()
	defer close(orderChan)

	for range unprocessedQty {
		if err := readAndProcessMessage(ctx, reader, orderChan); err != nil {
			slog.Error("Error processing order", "error", err)
		}
	}
}

func initReader() *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		GroupID: "orders-consumer-group",
		Topic:   "orders",
	})
}

func readAndProcessMessage(ctx context.Context, reader *kafka.Reader, orderChan chan models.Order) error {
	message, err := reader.ReadMessage(ctx)
	if err != nil {
		return fmt.Errorf("error reading order: %w", err)
	}

	if err := reader.CommitMessages(ctx, message); err != nil {
		return fmt.Errorf("error committing order: %w", err)
	}

	return orderToChannel(message.Value, orderChan)
}

func orderToChannel(message []byte, orderChan chan models.Order) error {
	var order models.Order

	if err := json.Unmarshal(message, &order); err != nil {
		return fmt.Errorf("error unmarshaling order: %w", err)
	}

	orderChan <- order
	slog.Info("Processed order", "order", order)
	return nil
}
