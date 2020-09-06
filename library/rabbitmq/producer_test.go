package rabbitmq

import (
	"fmt"
	"testing"
)

func Test_producer_Send(t *testing.T) {
	producer, _ := NewProducer()
	for i := 0; i < 10000; i++ {
		str := fmt.Sprintf("%d_hello_world开始发送消息测试", i + 1)
		producer.Send(TradeChange, str)
	}
	producer.Close()
}