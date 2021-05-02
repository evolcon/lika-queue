package memory

import (
	"github.com/lika_queue/brokers"
	"sync"
)

var lock sync.Mutex
var queues map[string]chan *brokers.Message

func init() {
	lock = sync.Mutex{}
	queues = make(map[string]chan *brokers.Message)
}

type Broker struct {
	Infinite bool
	len      int
}

func New(len int, infinite bool) *Broker {
	return &Broker{
		Infinite: infinite,
		len:      len,
	}
}

func (q *Broker) Publish(queue string, message interface{}, params map[string]interface{}) {
	q.PublishMessage(queue, &brokers.Message{Data: message, MQ: q}, params)
}

func (q *Broker) IsInfinite() bool {
	return q.Infinite
}

func (q *Broker) PublishMessage(queue string, message *brokers.Message, params map[string]interface{}) {
	message.QueueName = queue
	q.getOrCreateQueue(queue) <- message
}

func (q *Broker) Consume(queue string, params map[string]interface{}) *brokers.Message {
	mq := q.getOrCreateQueue(queue)

	if !q.Infinite && len(mq) == 0 {
		return nil
	}

	return <-mq
}

func (q *Broker) getOrCreateQueue(name string) chan *brokers.Message {
	lock.Lock()
	defer lock.Unlock()

	channel, ok := queues[name]

	if !ok {
		channel = make(chan *brokers.Message, q.len)
		queues[name] = channel
	}

	return channel
}
