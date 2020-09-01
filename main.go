package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	apiHTTP "gin-items/api/http"
	"gin-items/library/rabbitmq"
	"gin-items/library/setting"
)

func main() {

	router := apiHTTP.InitRouter()

	go initMqSubscriber()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.Config().Server.HttpPort),
		Handler:        router,
		ReadTimeout:    time.Duration(setting.Config().Server.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(setting.Config().Server.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}

func initMqSubscriber() {
	var (
		url = fmt.Sprintf("amqp://%s:%s@%s:%d/",
			setting.Config().RabbitMq.User,
			setting.Config().RabbitMq.Password,
			setting.Config().RabbitMq.Host,
			setting.Config().RabbitMq.Port)
		exchangeName = "service"
		exchangeKind = rabbitmq.ExchangeTopic
		routeKey     = "service_item"
		queueName    = "service.item"
	)

	mq, err := rabbitmq.NewMQ(url).Open()
	if err != nil {
		log.Printf("[ERROR] %s \n", err.Error())
		return
	}
	defer mq.Close()

	sub, err := mq.Subscriber("initSubscriber")
	if err != nil {
		log.Printf("[ERROR] Create subscriber failed, %v\n", err)
		return
	}
	defer sub.Close()

	exb := []*rabbitmq.ExchangeBinds{
		&rabbitmq.ExchangeBinds{
			Exch: rabbitmq.NewExchange(exchangeName, exchangeKind),
			Bindings: []*rabbitmq.Binding{
				&rabbitmq.Binding{
					RouteKey: routeKey,
					Queues: []*rabbitmq.Queue{
						rabbitmq.NewQueue(queueName),
					},
				},
			},
		},
	}

	msgCh := make(chan rabbitmq.Delivery, 1)
	defer close(msgCh)

	//sub.SetQos(10)
	if err = sub.SetExchangeBinds(exb).SetMagCallback(msgCh).SetQos(1).Open(); err != nil {
		log.Printf("[ERROR] Open failed, %v\n", err)
		return
	}

	for msg := range msgCh {
		log.Printf("Tag(%d) Body: %s\n", msg.DeliveryTag, string(msg.Body))
		msg.Ack(true)
	}
}
