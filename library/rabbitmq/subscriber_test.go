package rabbitmq

import (
	"fmt"
	"testing"
	"time"
)

func TestSubscriber(t *testing.T) {
	mq, err := NewMQ(mqUrl).Open()
	if err != nil {
		panic(err.Error())
	}
	defer mq.Close()

	sub, err := mq.Subscriber("test_subscriber")
	if err != nil {
		panic(fmt.Sprintf("Create subscriber failed, %v", err))
	}
	defer sub.Close()

	exb := []*ExchangeBinds{
		&ExchangeBinds{
			Exch: NewExchange("exch.unitest", ExchangeDirect),
			Bindings: []*Binding{
				&Binding{
					RouteKey: "route.unitest1",
					Queues: []*Queue{
						NewQueue("queue.unitest1"),
					},
				},
				&Binding{
					RouteKey: "route.unitest2",
					Queues: []*Queue{
						NewQueue("queue.unitest2"),
					},
				},
			},
		},
	}

	msgCh := make(chan Delivery, 1)
	defer close(msgCh)

	if err = sub.SetExchangeBinds(exb).SetMagCallback(msgCh).Open(); err != nil {
		panic(fmt.Sprintf("Open failed, %v", err))
	}

	i := 0
	for msg := range msgCh {
		i++
		if i%5 == 0 {
			sub.CloseChan()
		}
		t.Logf("Subscirber receive msg `%s`\n", string(msg.Body))
		time.Sleep(time.Second)
	}
}
