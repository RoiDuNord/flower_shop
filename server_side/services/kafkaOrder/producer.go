package kafkaorder

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"server/models"
	"sync"

	"github.com/segmentio/kafka-go"
)

func producer(ctx context.Context, reqQty int, writer *kafka.Writer) {
	var wg sync.WaitGroup

	wg.Add(reqQty)
	for i := 1; i <= reqQty; i++ {
		go func(i int) {
			defer wg.Done()
			if err := processOrder(ctx, i, writer); err != nil {
				slog.Error("Failed to load order", "orderID", i, "error", err)
			}
		}(i)
	}
	wg.Wait()
}

func processOrder(ctx context.Context, number int, writer *kafka.Writer) error {
	fileData, err := loadDataFromFile(number)
	if err != nil {
		return err
	}

	var order models.Order
	if err = json.Unmarshal(fileData, &order.BouquetsList); err != nil {
		return fmt.Errorf("error unmarshaling order data for order %d: %w", number, err)
	}

	if err = sendOrderToKafka(ctx, order, writer); err != nil {
		return fmt.Errorf("error sending order %d to Kafka: %w", number, err)
	}

	return nil
}

func sendOrderToKafka(ctx context.Context, order models.Order, writer *kafka.Writer) error {
	orderData, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order %d: %w", order.ID, err)
	}

	keyOrderID := fmt.Appendf(nil, "orderID_%d", order.ID)

	err = writer.WriteMessages(ctx,
		kafka.Message{
			Key:   keyOrderID,
			Value: orderData,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to write order %d to Kafka: %w", order.ID, err)
	}

	slog.Info("order sent to Kafka", "len", len(order.BouquetsList))
	return nil
}

func loadDataFromFile(number int) ([]byte, error) {
	file, err := getFilePath(number)
	if err != nil {
		return nil, fmt.Errorf("incorrect path %s", file)
	}
	paymentData, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("error reading from %s: %w", file, err)
	}
	return paymentData, nil
}

func getFilePath(orderNumber int) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %w", err)
	}
	return filepath.Join(cwd, "services", "kafkaOrder", "orders", fmt.Sprintf("order%d.json", orderNumber)), nil
}
