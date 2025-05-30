package kafkapayment

import (
	"context"
	"fmt"
	"log/slog"
	"server/config"
	kfk "server/services/initReaderWriter"
)

func SendPaymentToKafka(ctx context.Context, orderQty int) {
	cfg, err := config.GetKafkaParams("PAYMENTS")
	if err != nil {
		slog.Error("failed to load kafka configs", "error", err)
	}

	writer := kfk.InitWriter(cfg)
	defer writer.Close()

	producer(ctx, orderQty, writer)

	fmt.Println("Payment function exiting.")
}
