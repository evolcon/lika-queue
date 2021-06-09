package brokers

type BrokerInterface interface {
	// Publish Publishing messages
	Publish(queueName string, message interface{}, params map[string]interface{}) error
	// Consume Consuming messages
	Consume(queueName string, params map[string]interface{}) (MessageInterface, error)
}
