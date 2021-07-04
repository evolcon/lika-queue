package lika_queue

import (
	"sync"
)

type Broker struct {
	len    int
	lock   sync.Mutex
	queues map[string]chan MessageInterface
}

func NewInMemoryBroker(len int) BrokerInterface {
	return &Broker{
		len:    len,
		lock:   sync.Mutex{},
		queues: make(map[string]chan MessageInterface),
	}
}

func (q *Broker) Publish(queueName string, message interface{}, params map[string]interface{}) error {
	q.getOrCreateQueue(queueName) <- NewMessage(q, message, queueName, nil)

	return nil
}

func (q *Broker) Consume(queueName string, params map[string]interface{}) (MessageInterface, error) {
	mq := q.getOrCreateQueue(queueName)

	if len(mq) == 0 {
		return nil, nil
	}

	return <-mq, nil
}

func (q *Broker) getOrCreateQueue(name string) chan MessageInterface {
	q.lock.Lock()
	defer q.lock.Unlock()

	if channel, ok := q.queues[name]; ok {
		return channel
	}

	q.queues[name] = make(chan MessageInterface, q.len)

	return q.queues[name]
}
