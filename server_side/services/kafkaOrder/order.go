package kafkaorder

import (
	"context"
	"fmt"
	kfk "server/services/initReaderWriter"
)

func SendOrderToKafka(ctx context.Context, orderQty int) {
	writer := kfk.InitWriter("ORDERS")
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
