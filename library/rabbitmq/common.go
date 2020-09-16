package rabbitmq

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/streadway/amqp"
	"go.uber.org/zap"

	"gin-items/library/setting"
	"gin-items/middleware/log"
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
	url = fmt.Sprintf("amqp://%s:%s@%s:%d/",
		setting.Config().RabbitMq.User,
		setting.Config().RabbitMq.Password,
		setting.Config().RabbitMq.Host,
		setting.Config().RabbitMq.Port)
	exchangeName                           = "service"
	offLineReconnectInterval time.Duration = 10
	retryTimes                             = 5
)

func dealError(err error) error {
	if err != nil {
		log.ErrorLogger.Error("rabbitmq error", zap.String("error", err.Error()))
	}
	return err
}

func MqPack(data interface{}) ([]byte, error) {
	bytes, err := json.Marshal(data)
	return bytes, err
}

func MqUnpack(bytes []byte, data interface{}) error {
	err := json.Unmarshal(bytes, data)
	return err
}
