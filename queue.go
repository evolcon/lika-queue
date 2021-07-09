package lika_queue

import (
	"errors"
	"fmt"
)

type QueueInterface interface {
	// Add Add new one broker to pool
	Add(brokerName string, broker Broker)
	// DefaultBrokerName Returns default brokers name (key)
	DefaultBrokerName() string
	// MakeDefaultBroker Makes broker default by its name (key)
	MakeDefaultBroker(brokerName string) (e error)
	// Publish Sending message to default broker
	Publish(queueName string, message interface{}, params map[string]interface{}) error
	// Consume Consuming messages from default broker
	Consume(queueName string, params map[string]interface{}) (MessageData, error)
	// DefaultBroker Returns default broker
	DefaultBroker() (Broker, error)
	// Broker Returns broker by the name
	Broker(brokerName string) (b Broker, e error)
}

func New() QueueInterface {
	return &Queue{brokers: make(map[string]Broker)}
}

type Queue struct {
	defaultBrokerName string
	brokers           map[string]Broker
}

func (q *Queue) Add(brokerName string, broker Broker) {
	q.brokers[brokerName] = broker

	if q.defaultBrokerName == "" {
		q.defaultBrokerName = brokerName
	}
}

func (q *Queue) DefaultBrokerName() string {
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
	broker, err := q.DefaultBroker()

	if err != nil {
		panic("Default broker does not set")
	}

	return broker.Publish(queueName, message, params)
}

func (q *Queue) Consume(queueName string, params map[string]interface{}) (MessageData, error) {
	broker, err := q.DefaultBroker()

	if err != nil {
		panic("Default broker does not set")
	}

	return broker.Consume(queueName, params)
}

func (q *Queue) DefaultBroker() (Broker, error) {
	return q.Broker(q.defaultBrokerName)
}

func (q *Queue) Broker(brokerName string) (b Broker, e error) {
	b, ok := q.brokers[brokerName]

	if !ok {
		e = fmt.Errorf("broker '%v' does not exist", brokerName)
	}

	return
}
