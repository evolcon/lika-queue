package lika_queue

import (
	"errors"
	"fmt"
	"github.com/lika_queue/brokers"
)

type QueueInterface interface {
	Add(name string, broker brokers.BrokerInterface)
	GetDefaultBrokerName() string
	MakeDefaultBroker(brokerName string) (e error)
	Publish(queueName string, message interface{}, params map[string]interface{}) error
	Consume(queueName string, params map[string]interface{}) (brokers.MessageInterface, error)
	DefaultBroker() (brokers.BrokerInterface, error)
	Broker(brokerName string) (b brokers.BrokerInterface, e error)
}

func New() QueueInterface {
	return &Queue{brokers: make(map[string]brokers.BrokerInterface)}
}

type Queue struct {
	defaultBrokerName string
	brokers           map[string]brokers.BrokerInterface
}

func (q *Queue) Add(name string, broker brokers.BrokerInterface) {
	q.brokers[name] = broker

	if q.defaultBrokerName == "" {
		q.defaultBrokerName = name
	}
}

func (q *Queue) GetDefaultBrokerName() string {
	return q.defaultBrokerName
}

func (q *Queue) MakeDefaultBroker(brokerName string) (e error) {
	if q.brokers[brokerName] == nil {
		return errors.New(fmt.Sprintf("Message Broker Broker named '%v' does not exist", brokerName))
	}

	q.defaultBrokerName = brokerName

	return
}

func (q *Queue) Publish(queueName string, message interface{}, params map[string]interface{}) error {
	connection, err := q.DefaultBroker()

	if err != nil {
		panic("Default connection does not set")
	}

	return connection.Publish(queueName, message, params)
}

func (q *Queue) Consume(queueName string, params map[string]interface{}) (brokers.MessageInterface, error) {
	connection, err := q.DefaultBroker()

	if err != nil {
		panic("Default connection does not set")
	}

	return connection.Consume(queueName, params)
}

func (q *Queue) DefaultBroker() (brokers.BrokerInterface, error) {
	return q.Broker(q.defaultBrokerName)
}

func (q *Queue) Broker(brokerName string) (b brokers.BrokerInterface, e error) {
	b, ok := q.brokers[brokerName]

	if !ok {
		e = fmt.Errorf("broker '%v' does not exist", brokerName)
	}

	return
}
