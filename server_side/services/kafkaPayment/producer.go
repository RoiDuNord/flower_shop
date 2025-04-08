package kafkapayment

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
			if err := paymentToKafka(ctx, i, writer); err != nil {
				slog.Error("Failed to load payment", "paymentID", i, "error", err)
			}
		}(i)
	}
	wg.Wait()
}

func paymentToKafka(ctx context.Context, number int, producer *kafka.Writer) error {
	fileData, err := loadDataFromFile(number)
	if err != nil {
		return err
	}

	var payment models.Payment
	if err = json.Unmarshal(fileData, &payment); err != nil {
		return fmt.Errorf("error unmarshaling payment data: %w", err)
	}
	fmt.Println("payment", payment.OrderID)

	if err = sendPaymentToKafka(ctx, payment, producer); err != nil {
		return fmt.Errorf("error sending payment to Kafka: %w", err)
	}

	return nil
}

func sendPaymentToKafka(ctx context.Context, payment models.Payment, producer *kafka.Writer) error {
	paymentData, err := json.Marshal(payment)
	if err != nil {
		return fmt.Errorf("failed to marshal payment: %w", err)
	}

	err = producer.WriteMessages(ctx,
		kafka.Message{
			Key:   fmt.Appendf(nil, "paymentID_%d", payment.PaymentID),
			Value: paymentData,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to write payment to Kafka: %w", err)
	}

	// slog.Info("payment sent to Kafka", "paymentID", payment.PaymentID)
	return nil
}

func loadDataFromFile(number int) ([]byte, error) {
	file := getFilePath(number)
	paymentData, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("error reading from %s: %w", file, err)
	}
	return paymentData, nil
}

func getFilePath(paymentNumber int) string {
	cwd, err := os.Getwd()
	if err != nil {
		slog.Error("Error getting current directory", "error", err)
		return ""
	}
	return filepath.Join(cwd, "services", "kafkaPayment", "payments", fmt.Sprintf("payment%d.json", paymentNumber))
}
