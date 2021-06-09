package memory

import (
	"github.com/lika_queue/brokers"
	"sync"
)

var lock sync.Mutex
var queues map[string]chan brokers.MessageInterface

func init() {
	lock = sync.Mutex{}
	queues = make(map[string]chan brokers.MessageInterface)
}

type Broker struct {
	len int
}

func New(len int) brokers.BrokerInterface {
	return &Broker{
		len: len,
	}
}

func (q *Broker) Publish(queue string, message interface{}, params map[string]interface{}) error {
	q.getOrCreateQueue(queue) <- brokers.NewMessage(q, message, queue, nil)

	return nil
}

func (q *Broker) Consume(queue string, params map[string]interface{}) (brokers.MessageInterface, error) {
	mq := q.getOrCreateQueue(queue)

	if len(mq) == 0 {
		return nil, nil
	}

	return <-mq, nil
}

func (q *Broker) getOrCreateQueue(name string) chan brokers.MessageInterface {
	lock.Lock()
	defer lock.Unlock()

	channel, ok := queues[name]

	if !ok {
		channel = make(chan brokers.MessageInterface, q.len)
		queues[name] = channel
	}

	return channel
}
