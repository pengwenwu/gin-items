package rabbitmq

import (
	"github.com/streadway/amqp"
	"time"
)

// 交换机类型 Exchange type
var (
	ExchangeDirect  = amqp.ExchangeDirect
	ExchangeFanout  = amqp.ExchangeFanout
	ExchangeTopic   = amqp.ExchangeTopic
	ExchangeHeaders = amqp.ExchangeHeaders
)

// 持久化类型 DeliveryMode
var (
	Transient  uint8 = amqp.Transient
	Persistent uint8 = amqp.Persistent
)

// ExchangeBinds exchange ==> routeKey ==> queues
type ExchangeBinds struct {
	Exch     *Exchange
	Bindings []*Binding
}

// Binding routeKey ==> queues
type Binding struct {
	RouteKey string
	Queues   []*Queue
	NoWait   bool
	Args     amqp.Table
}

// 交换机配置 Exchange
type Exchange struct {
	Name       string
	Kind       string // 类型
	Durable    bool   // 持久化
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

func NewExchange(name string, kind string) *Exchange {
	return &Exchange{
		Name:       name,
		Kind:       kind,
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
		Args:       nil,
	}
}

// Queue 配置
type Queue struct {
	Name       string
	Durable    bool // 持久化
	AutoDelete bool
	Exclusive  bool // 排他
	NoWait     bool
	Args       amqp.Table
}

func NewQueue(name string) *Queue {
	return &Queue{
		Name:       name,
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
		Args:       nil,
	}
}

// 生产者数据格式
type PublisherMsg struct {
	ContentType     string // MIME content type
	ContentEncoding string // MIME content type
	DeliverMode     uint8  // Transient or Persistent
	Priority        uint8  // 0 to 9
	Timestamp       time.Time
	Body            []byte
}

func NewPublishMsg(body []byte) *PublisherMsg {
	return &PublisherMsg{
		ContentType:     "application/json",
		ContentEncoding: "",
		DeliverMode:     Persistent,
		Priority:        uint8(5),
		Timestamp:       time.Now(),
		Body:            body,
	}
}

// 消费者消费选项
type SubscribeOption struct {
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
}

func DefaultSubscribeOption() *SubscribeOption {
	return &SubscribeOption{
		NoWait: true,
	}
}
