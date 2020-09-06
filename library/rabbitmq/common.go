package rabbitmq

import (
	"time"

	"github.com/streadway/amqp"
)

// Exchange type
var (
	ExchangeDirect  = amqp.ExchangeDirect
	ExchangeFanout  = amqp.ExchangeFanout
	ExchangeTopic   = amqp.ExchangeTopic
	ExchangeHeaders = amqp.ExchangeHeaders
)

// DeliveryMode
var (
	Transient  uint8 = amqp.Transient
	Persistent uint8 = amqp.Persistent
)

var (
	url                                    = "amqp://guest:guest@localhost:5672/"
	exchangeName                           = "service"
	offLineReconnectInterval time.Duration = 10
	retryTimes                             = 5
)

// @todo 定义一个错误处理函数
func dealError(err error) error {
	if err != nil {
		//variable.ZapLog.Error(err.Error())
	}
	return err
}
