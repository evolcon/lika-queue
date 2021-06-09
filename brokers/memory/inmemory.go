package memory

import (
	queue "github.com/lika_queue"
	"sync"
)

type Broker struct {
	len    int
	lock   sync.Mutex
	queues map[string]chan queue.MessageInterface
}

func New(len int) queue.BrokerInterface {
	return &Broker{
		len:    len,
		lock:   sync.Mutex{},
		queues: make(map[string]chan queue.MessageInterface),
	}
}

func (q *Broker) Publish(queueName string, message interface{}, params map[string]interface{}) error {
	q.getOrCreateQueue(queueName) <- queue.NewMessage(q, message, queueName, nil)

	return nil
}

func (q *Broker) Consume(queueName string, params map[string]interface{}) (queue.MessageInterface, error) {
	mq := q.getOrCreateQueue(queueName)

	if len(mq) == 0 {
		return nil, nil
	}

	return <-mq, nil
}

func (q *Broker) getOrCreateQueue(name string) chan queue.MessageInterface {
	q.lock.Lock()
	defer q.lock.Unlock()

	channel, ok := q.queues[name]

	if !ok {
		channel = make(chan queue.MessageInterface, q.len)
		q.queues[name] = channel
	}

	return channel
}
