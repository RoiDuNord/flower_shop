package kafkaorder

import (
	"context"
	"log/slog"
	"server/config"
	kfk "server/services/initReaderWriter"
)

func SendOrderToKafka(ctx context.Context, orderQty int) {
	cfg, err := config.GetKafkaParams("ORDERS")
	if err != nil {
		slog.Error("failed to load kafka configs", "error", err)
	}

	writer := kfk.InitWriter(cfg)
	defer writer.Close()

	producer(ctx, orderQty, writer)

	slog.Info("SendOrderToKafka function exiting")
}
