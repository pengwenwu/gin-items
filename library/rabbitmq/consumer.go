package rabbitmq

import (
	"time"

	"github.com/streadway/amqp"
)

type consumer struct {
	conn                        *amqp.Connection
	exchangeType                string
	exchangeName                string
	queueBind                   *queueBind
	durable                     bool
	occurErr                    error // 记录初始化过程中的错误
	connErr                     chan *amqp.Error
	routeKey                    string                   // 断线重连，结构体内部使用
	callbackForReceived         func(receiveData []byte) // 断线重连，结构体内部使用
	offlineReconnectIntervalSec time.Duration
	retryTimes                  int
	callbackOffline             func(err *amqp.Error)
}

func NewConsumer() (*consumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	consumer := &consumer{
		conn:                        conn,
		exchangeType:                ExchangeTopic,
		exchangeName:                exchangeName,
		durable:                     true,
		connErr:                     conn.NotifyClose(make(chan *amqp.Error, 1)),
		callbackForReceived:         nil,
		offlineReconnectIntervalSec: offLineReconnectInterval,
		retryTimes:                  retryTimes,
	}
	return consumer, nil
}

func (c *consumer) Received(queueBind *queueBind, callbackFuncDealMsg func(receivedData []byte)) {
	defer func() {
		_ = c.conn.Close()
	}()

	// 将回调函数地址赋值给结构体变量，用于掉线重连使用
	c.queueBind = queueBind
	c.callbackForReceived = callbackFuncDealMsg

	blocking := make(chan bool)

	go func() {
		ch, err := c.conn.Channel()
		c.occurErr = dealError(err)
		defer ch.Close()

		// 声明交换机
		err = ch.ExchangeDeclare(
			c.exchangeName,
			c.exchangeType,
			true,
			false,
			false,
			false,
			nil,
		)

		// 声明队列
		queue, err := ch.QueueDeclare(
			c.queueBind.name,
			true,
			false,
			false,
			false,
			nil,
		)
		c.occurErr = dealError(err)

		// 队列绑定
		for _, routeKey := range c.queueBind.keys {
			err = ch.QueueBind(
				queue.Name,
				string(routeKey),
				c.exchangeName,
				false,
				nil,
			)
			c.occurErr = dealError(err)
		}

		msgs, err := ch.Consume(
			c.queueBind.name,
			"",
			true,
			false,
			false,
			false,
			nil,
		)

		for msg := range msgs {
			callbackFuncDealMsg(msg.Body)
		}
	}()

	<-blocking
}
