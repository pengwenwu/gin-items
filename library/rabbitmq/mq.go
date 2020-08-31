package rabbitmq

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

var (
	StateClosed    = uint8(0)
	StateOpened    = uint8(1)
	StateReopening = uint8(2)
)

type MQ struct {
	url         string           // 连接url
	mutex       sync.RWMutex     // 读写锁
	conn        *amqp.Connection // RabbitMq tcp 连接
	publishers  []*Publisher
	subscribers []*Subscriber
	closeCh     chan *amqp.Error // RabbitMq 监听连接错误
	stopCh      chan struct{}    // 监听用户手动关闭
	state       uint8            // MQ状态
}

func NewMQ(url string) *MQ {
	return &MQ{
		url:         url,
		mutex:       sync.RWMutex{},
		publishers:  make([]*Publisher, 0, 1),
		state:       StateClosed,
	}
}

func (m *MQ) Open() (mq *MQ, err error) {
	// 进行open期间不允许做任何跟连接有关的事情
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.state == StateOpened {
		return m, errors.New("MQ: Had been opened")
	}
	if m.conn, err = m.dial(); err != nil {
		return m, fmt.Errorf("MQ: Dial err: %v", err)
	}

	m.state = StateOpened
	m.stopCh = make(chan struct{})
	m.closeCh = make(chan *amqp.Error, 1)
	m.conn.NotifyClose(m.closeCh)

	go m.keepalive()

	return m, nil
}

func (m *MQ) Close() {
	m.mutex.Lock()

	// 关闭所有生产者 close publisher
	for _, p := range m.publishers {
		p.Close()
	}
	m.publishers = m.publishers[:0]

	// close subscriber
	for _, s := range m.subscribers {
		s.Close()
	}
	m.subscribers = m.subscribers[:0]

	// close mq connection
	select {
	case <-m.stopCh:
		// had been closed
	default:
		close(m.stopCh)
	}

	m.mutex.Unlock()

	// wait done
	for m.State() != StateClosed {
		time.Sleep(time.Second)
	}
}

func (m *MQ) Publisher(name string) (*Publisher, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.state != StateOpened {
		return nil, fmt.Errorf("MQ: Not initialized, now state is %d", m.state)
	}

	p := newPublisher(name, m)
	m.publishers = append(m.publishers, p)
	return  p,nil
}

func (m *MQ) Subscriber(name string) (*Subscriber, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.state != StateOpened {
		return nil, fmt.Errorf("MQ: Not initialized, now state is %d", m.state)
	}
	s := newSubscriber(name, m)
	m.subscribers = append(m.subscribers, s)
	return s, nil
}

func (m *MQ) State() uint8 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.state
}

func (m *MQ) keepalive() {
	select {
	case <-m.stopCh:
		// 正常关闭
		log.Println("[WARN] MQ: Shutdown normally.")
		m.mutex.Lock()
		m.conn.Close()
		m.state = StateClosed
		m.mutex.Unlock()
	case err := <-m.closeCh:
		if err == nil {
			log.Println("[ERROR] MQ: Disconnected with MQ, but Error detail is nil")
		} else {
			log.Printf("[ERROR] MQ: Disconnected whit MQ, code: %d, reason: %s\n", err.Code, err.Reason)
		}

		// tcp连接中断，重新连接
		m.mutex.Lock()
		m.state = StateReopening
		m.mutex.Unlock()

		maxRetry := 9999999
		for i := 0; i < maxRetry; i++ {
			time.Sleep(time.Second)
			if _, err := m.Open(); err != nil {
				log.Printf("[ERROR] MQ: Connection recover failed for %d times, %v\n", i+1, err)
				continue
			}
			log.Printf("[INFO] MQ: Connection recover OK. Total try %d times\n", i+1)
			return
		}
		log.Printf("[ERROR] MQ: Try to reconnect to MQ failed over maxRetry(%d), so exit.\n", maxRetry)
	}
}

func (m *MQ) channel() (*amqp.Channel, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.conn.Channel()
}

func (m *MQ) dial() (*amqp.Connection, error) {
	return amqp.Dial(m.url)
}
