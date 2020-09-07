package rabbitmq

import "github.com/streadway/amqp"

type producer struct {
	conn         *amqp.Connection
	exchangeType string
	exchangeName string
	occurErr     error
}

func NewProducer() (*producer, error) {
	conn, err := amqp.Dial(url)

	if err != nil {
		dealError(err)
		return nil, err
	}

	producer := &producer{
		conn:         conn,
		exchangeType: ExchangeTopic,
		exchangeName: exchangeName,
		occurErr:     nil,
	}
	return producer, nil
}

func (p *producer) Send(routeKey RouteKey, data []byte) bool {
	ch, err := p.conn.Channel()
	p.occurErr = dealError(err)
	defer ch.Close()

	// 声明交换机，该模式生产者只负责将消息投递到交换机即可
	err = ch.ExchangeDeclare(
		p.exchangeName,
		p.exchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
	p.occurErr = dealError(err)

	// 投递消息
	err = ch.Publish(
		p.exchangeName,
		string(routeKey),
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: Persistent,
			Body:         data,
		},
	)
	p.occurErr = dealError(err)
	if p.occurErr != nil {
		return false
	}
	return true
}

// 发送完毕手动关闭，这样不影响send多次发送数据
func (p *producer) Close() {
	p.conn.Close()
}
