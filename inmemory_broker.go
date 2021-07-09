package lika_queue

import (
	"sync"
)

type MemoryBroker struct {
	len    int
	lock   sync.Mutex
	queues map[string]chan MessageData
}

func NewMemoryBroker(len int) Broker {
	return &MemoryBroker{
		len:    len,
		lock:   sync.Mutex{},
		queues: make(map[string]chan MessageData),
	}
}

func (q *MemoryBroker) Publish(queueName string, message interface{}, params map[string]interface{}) error {
	q.getOrCreateQueue(queueName) <- NewMessage(q, message, queueName, nil)

	return nil
}

func (q *MemoryBroker) Consume(queueName string, params map[string]interface{}) (MessageData, error) {
	mq := q.getOrCreateQueue(queueName)

	if len(mq) == 0 {
		return nil, nil
	}

	return <-mq, nil
}

func (q *MemoryBroker) getOrCreateQueue(name string) chan MessageData {
	q.lock.Lock()
	defer q.lock.Unlock()

	if channel, ok := q.queues[name]; ok {
		return channel
	}

	q.queues[name] = make(chan MessageData, q.len)

	return q.queues[name]
}
