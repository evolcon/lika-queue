package brokers

type BrokerInterface interface {
	Consume(queueName string, params map[string]interface{}) *Message
	Publish(queueName string, message interface{}, params map[string]interface{})
	PublishMessage(queueName string, message *Message, params map[string]interface{})
	IsInfinite() bool
}
