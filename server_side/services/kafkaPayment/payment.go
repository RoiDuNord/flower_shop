package kafkapayment

import (
	"context"
	"fmt"
	kfk "server/services/initReaderWriter"
)

func SendPaymentToKafka(ctx context.Context, orderQty int) {
	writer := kfk.InitWriter("PAYMENTS")
	defer writer.Close()

	producer(ctx, orderQty, writer)

	fmt.Println("Payment function exiting.")
}
