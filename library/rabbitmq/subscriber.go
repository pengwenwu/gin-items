package rabbitmq

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

type Delivery struct {
	amqp.Delivery
}

type Subscriber struct {
	name          string
	mq            *MQ              // MQ实例
	mutex         sync.RWMutex    // 读写锁
	ch            *amqp.Channel    // MQ的会话channel
	exchangeBinds []*ExchangeBinds // MQ的exchange与其绑定的queues
	prefetch      int              // Qos prefetch
	callback      chan<- Delivery  // 上层用于接收消费出来的消息的管道
	closeCh       chan *amqp.Error
	stopCh        chan struct{}
	state         uint8
}

func newSubscriber(name string, mq *MQ) *Subscriber {
	return &Subscriber{
		name:   name,
		mq:     mq,
		stopCh: make(chan struct{}),
	}
}

func (s *Subscriber) Name() string {
	return s.name
}

// CloseChan 该接口仅用于测试使用，勿手动调用
func (s *Subscriber) CloseChan() {
	s.mutex.Lock()
	s.ch.Close()
	s.mutex.Unlock()
}

func (s *Subscriber) SetExchangeBinds(eb []*ExchangeBinds) *Subscriber {
	s.mutex.Lock()
	if s.state != StateOpened {
		s.exchangeBinds = eb
	}
	s.mutex.Unlock()
	return s
}

func (s *Subscriber) SetMagCallback(cb chan<- Delivery) *Subscriber {
	s.mutex.Lock()
	s.callback = cb
	s.mutex.Unlock()
	return s
}

// SetQos 设置channel粒度的Qos，prefetch取值范围[0, ∞]，默认为0
// 如果想要RoundRobin轮流地进行消费，设置prefetch为1即可
// 注意：在调用open之前设置
func (s *Subscriber) SetQos(prefetch int) *Subscriber {
	s.mutex.Lock()
	s.prefetch = prefetch
	s.mutex.Unlock()
	return s
}

func (s *Subscriber) Open() error {
	// Open期间不允许对channel做任何操作
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 参数校验
	if s.mq == nil {
		return errors.New("MQ: Bad subscriber")
	}
	if len(s.exchangeBinds) <= 0 {
		return errors.New("MQ: No exchangeBinds found. You should SetExchangeBinds before")
	}

	// 状态检测
	if s.state == StateOpened {
		return errors.New("MQ: Subscriber had been opened")
	}

	// 初始化channel
	ch, err := s.mq.channel()
	if err != nil {
		return fmt.Errorf("MQ: Create channel failed, %v", err)
	}

	err = func(ch *amqp.Channel) error {
		var e error
		if e = applyExchangeBinds(ch, s.exchangeBinds); e != nil {
			return e
		}
		if e = ch.Qos(s.prefetch, 0, false); e != nil {
			return e
		}
		return nil
	}(ch)
	if err != nil {
		ch.Close()
		return fmt.Errorf("MQ: %v", err)
	}

	s.ch = ch
	s.state = StateOpened
	s.stopCh = make(chan struct{})
	s.closeCh = make(chan *amqp.Error, 1)
	s.ch.NotifyClose(s.closeCh)

	// 开始循环消费
	opt := DefaultSubscribeOption()
	notify := make(chan error, 1)
	s.subscribe(opt, notify)
	for e := range notify {
		if e != nil {
			log.Printf("[ERROR] %v\n", e)
			continue
		}
	}
	close(notify)

	// 健康检测
	go s.keepalive()

	return nil
}


func (s *Subscriber) Close() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	select {
	case <-s.stopCh:
		// had been closed
	default:
		close(s.stopCh)
	}
}

// notifyErr 向上层抛出错误, 如果error为空表示执行完成.由上层负责关闭channel
func (s *Subscriber) subscribe(opt *SubscribeOption, notifyErr chan<- error) {
	for idx, eb := range s.exchangeBinds {
		if eb == nil {
			notifyErr <- fmt.Errorf("MQ: ExchangeBinds[%d] is nil, subscriber(%s)", idx, s.name)
			continue
		}
		for i, b := range eb.Bindings {
			if b == nil {
				notifyErr <- fmt.Errorf("MQ: Binding[%d] is nil, ExchangeBinds[%d], consumer(%s)", i, idx, s.name)
				continue
			}
			for qi, q := range b.Queues {
				if q == nil {
					notifyErr <- fmt.Errorf("MQ: Queue[%d] is nil, ExchangeBinds[%d], Biding[%d], consumer(%s)", qi, idx, i, s.name)
					continue
				}
				delivery, err := s.ch.Consume(q.Name, "", opt.AutoAck, opt.Exclusive, opt.NoLocal, opt.NoWait, opt.Args)
				if err != nil {
					notifyErr <- fmt.Errorf("MQ: Consumer(%s) consume queue(%s) failed, %v", s.name, q.Name, err)
					continue
				}
				go s.deliver(delivery)
			}
		}
	}
	notifyErr <- nil
}

func (s *Subscriber) deliver(delivery <-chan amqp.Delivery) {
	for d := range delivery {
		if s.callback != nil {
			s.callback <- Delivery{d}
		}
	}
}

func (s *Subscriber) State() uint8 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.state
}

func (s *Subscriber) keepalive() {
	select {
	case <-s.stopCh:
		// 正常关闭
		log.Printf("[WARN] Subscriber(%s) shutdown normally\n", s.Name())
		s.mutex.Lock()
		s.ch.Close()
		s.ch = nil
		s.state = StateClosed
		s.mutex.Unlock()
	case err := <-s.closeCh:
		if err == nil {
			log.Printf("[ERROR] MQ: Subscirber(%s)'s channel was closed, but Error detail is nil\n", s.name)
		} else {
			log.Printf("[ERROR] MQ: Subscriber(%s)'s channel was cloed, code: %d, reason: %s\n", s.name, err.Code, err.Reason)
		}

		maxRetry := 99999999
		for i := 0; i < maxRetry; i++ {
			time.Sleep(time.Second)
			if s.mq.State() != StateOpened {
				log.Printf("[WARN] MQ: Consumer(%s) try to recover channel for %d times, but mq's state != StateOpened\n", s.name, i+1)
				continue
			}
			if e := s.Open(); e != nil {
				log.Printf("[WARN] MQ: Consumer(%s) recover channel failed for %d times, Err:%v\n", s.name, i+1, e)
				continue
			}
			log.Printf("[INFO] MQ: Consumer(%s) recover channel OK. Total try %d times\n", s.name, i+1)
			return
		}
		log.Printf("[ERROR] MQ: Consumer(%s) try to recover channel over maxRetry(%d), so exit\n", s.name, maxRetry)
	}
}