package kafkaorder

import (
	"context"
	"fmt"
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

	fmt.Println("Order function exiting.")
}

// go consumer(ctx, orderCh, orderQty)

// printCh(orderCh)

// func printCh(orderCh chan models.Order) {
// 	for order := range orderCh {
// 		fmt.Println("Received order in main:", order)
// 	}
// }
