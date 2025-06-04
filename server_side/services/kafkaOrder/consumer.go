package kafkaorder

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"server/config"
	"server/models"
	kfk "server/services/initReaderWriter"

	"github.com/segmentio/kafka-go"
)

func Consumer(ctx context.Context, cfg config.KafkaParams, orderChan chan models.Order, unprocessedQty int) {
	reader := kfk.InitReader(cfg)
	slog.Info("kafka-orders is ready for consuming")
	defer reader.Close()
	defer close(orderChan)

	for range unprocessedQty {
		if err := consumeAndProcessMessage(ctx, reader, orderChan); err != nil {
			slog.Error("Error processing order", "error", err)
			continue
		}
	}
}

func consumeAndProcessMessage(ctx context.Context, reader *kafka.Reader, orderChan chan models.Order) error {
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
