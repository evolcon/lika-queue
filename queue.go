package lika_queue

import (
	"errors"
	"fmt"
	"github.com/lika_queue/brokers"
)

func New() *Queue {
	return &Queue{brokers: make(map[string]brokers.BrokerInterface)}
}

type Queue struct {
	defaultBrokerName string
	brokers           map[string]brokers.BrokerInterface
}

func (q *Queue) Add(name string, connection brokers.BrokerInterface) {
	q.brokers[name] = connection

	if q.defaultBrokerName == "" {
		q.defaultBrokerName = name
	}
}

func (q *Queue) GetDefaultBrokerName() string {
	return q.defaultBrokerName
}

func (q *Queue) MakeDefaultBroker(name string) (e error) {
	if q.brokers[name] == nil {
		return errors.New(fmt.Sprintf("Message Queue Broker named '%v' does not exist", name))
	}

	q.defaultBrokerName = name

	return
}

func (q *Queue) Publish(queueName string, message interface{}, params map[string]interface{}) {
	connection, err := q.GetDefaultBroker()

	if err != nil {
		panic("Default connection does not set")
	}

	connection.Publish(queueName, message, params)
}

func (q *Queue) Consume(queueName string, params map[string]interface{}) *brokers.Message {
	connection, err := q.GetDefaultBroker()

	if err != nil {
		panic("Default connection does not set")
	}

	return connection.Consume(queueName, params)
}

func (q *Queue) GetDefaultBroker() (brokers.BrokerInterface, error) {
	return q.brokers[q.defaultBrokerName], q.checkConnection(q.defaultBrokerName)
}

func (q *Queue) checkConnection(name string) (e error) {
	if q.brokers[name] == nil {
		e = errors.New(fmt.Sprintf("Connection named '%v' does not exist", name))
	}

	return
}
