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
	Infinite bool
	len      int
}

func New(len int, infinite bool) brokers.BrokerInterface {
	return &Broker{
		Infinite: infinite,
		len:      len,
	}
}

func (q *Broker) Publish(queue string, message interface{}, params map[string]interface{}) error {
	return q.PublishMessage(queue, brokers.NewMessage(q, message, queue, make(map[string]interface{})), params)
}

func (q *Broker) IsInfinite() bool {
	return q.Infinite
}

func (q *Broker) PublishMessage(queue string, message brokers.MessageInterface, params map[string]interface{}) error {
	q.getOrCreateQueue(queue) <- brokers.NewMessage(q, message.GetData(), queue, message.GetMetaData())

	return nil
}

func (q *Broker) Consume(queue string, params map[string]interface{}) (brokers.MessageInterface, error) {
	mq := q.getOrCreateQueue(queue)

	if !q.Infinite && len(mq) == 0 {
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
