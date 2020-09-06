package rabbitmq

import (
	"fmt"
	"testing"
)

func Test_consumer_Received(t *testing.T) {
	consumer, _ := NewConsumer()
	consumer.Received(OrderUserRelCreateUpdate, func(receivedData string) {
		fmt.Printf("topic回调函数处理消息：--->%s\n", receivedData)
	})
}