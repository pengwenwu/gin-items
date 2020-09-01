package rabbitmq

import (
	"fmt"
	"testing"
	"time"
)

const mqUrl = "amqp://guest:guest@localhost:5672/"

func TestPublisher(t *testing.T) {
	mq, err := NewMQ(mqUrl).Open()
	if err != nil {
		panic(err.Error())
	}

	pub, err := mq.Publisher("test-publisher")
	if err != nil {
		panic(fmt.Sprintf("Create publisher failed, %v", err))
	}

	exb := []*ExchangeBinds{
		&ExchangeBinds{
			Exch:     NewExchange("exch.unitest", ExchangeDirect),
			Bindings: []*Binding{
				&Binding{
					RouteKey: "route.unitest1",
					Queues:   []*Queue{
						NewQueue("queue.unitest1"),
					},
				},
				&Binding{
					RouteKey: "route.unitest2",
					Queues:   []*Queue{
						NewQueue("queue.unitest2"),
					},
				},
			},
		},
	}

	if err = pub.SetExchangeBinds(exb).Confirm(true).Open(); err != nil {
		panic(fmt.Sprintf("Open failed, %v", err))
	}

	for i := 0; i < 4; i++ {
		if i > 0 && i%3 == 0 {
			pub.CloseChan()
		}
		err = pub.Publish("exch.unitest", "route.unitest2", NewPublishMsg([]byte(`{"name":"pww"}`)))
		t.Logf("Publisher state: %d, err: %v\n", pub.State(), err)
		time.Sleep(time.Second)
	}

	pub.Close()
	mq.Close()
}