package brokers

type BrokerInterface interface {
	Consume(queueName string, params map[string]interface{}) (MessageInterface, error)
	Publish(queueName string, message interface{}, params map[string]interface{}) error
	PublishMessage(queueName string, message MessageInterface, params map[string]interface{}) error
	IsInfinite() bool
}
