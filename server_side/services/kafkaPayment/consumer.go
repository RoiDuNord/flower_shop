package kafkapayment

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"server/models"
	"time"

	"github.com/segmentio/kafka-go"
)

func Consumer(ctx context.Context, paymentChan chan models.Payment, unprocessedQty int) {
	t := time.Now()
	reader := initReader()
	defer reader.Close()
	defer close(paymentChan)

	for range unprocessedQty {
		if err := readAndProcessMessage(ctx, reader, paymentChan); err != nil {
			slog.Error("Error processing order", "error", err)
		}
	}

	// for range unprocessedQty {
	// 	go func() {
	// 		if err := readAndProcessMessage(ctx, reader, paymentChan); err != nil {
	// 			slog.Error("Error processing order", "error", err)
	// 		}
	// 	}()
	// }

	dur := time.Since(t)
	fmt.Println("Payment", dur)
}

func initReader() *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		GroupID: "payments-consumer-group",
		Topic:   "payments",
	})
}

func readAndProcessMessage(ctx context.Context, reader *kafka.Reader, paymentChan chan models.Payment) error {
	message, err := reader.ReadMessage(ctx)
	if err != nil {
		return fmt.Errorf("error reading payment: %w", err)
	}

	if err := reader.CommitMessages(ctx, message); err != nil {
		return fmt.Errorf("error committing payment: %w", err)
	}

	return paymentToChannel(message.Value, paymentChan)
}

func paymentToChannel(message []byte, paymentChan chan models.Payment) error {
	var payment models.Payment
	if err := json.Unmarshal(message, &payment); err != nil {
		return fmt.Errorf("error unmarshaling payment: %w", err)
	}

	paymentChan <- payment
	slog.Info("Processed payment", "payment", payment)
	return nil
}
