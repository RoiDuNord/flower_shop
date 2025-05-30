package kafkapayment

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"server/config"
	"server/models"
	kfk "server/services/initReaderWriter"
	"time"

	"github.com/segmentio/kafka-go"
)

func Consumer(ctx context.Context, cfg config.KafkaParams, paymentChan chan models.Payment, unprocessedQty int) {
	t := time.Now()

	reader := kfk.InitReader(cfg)
	slog.Info("kafka-payments is ready for consuming")
	defer reader.Close()
	defer close(paymentChan)

	for range unprocessedQty {
		if err := consumeAndProcessMessage(ctx, reader, paymentChan); err != nil {
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

func consumeAndProcessMessage(ctx context.Context, reader *kafka.Reader, paymentChan chan models.Payment) error {
	message, err := reader.ReadMessage(ctx)
	if err != nil {
		return fmt.Errorf("error reading payment: %w", err)
	}
	fmt.Println("message", string(message.Value))

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
