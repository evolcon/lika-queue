package lika_queue

type MessageData interface {
	Data() interface{}
	MetaData() map[string]interface{}
	QueueName() string
}

func NewMessage(broker Broker, data interface{}, queueName string, metaData map[string]interface{}) MessageData {
	return &Message{
		data:      data,
		broker:    broker,
		queueName: queueName,
		metaData:  metaData,
	}
}

type Message struct {
	metaData  map[string]interface{}
	data      interface{}
	retries   int
	queueName string
	done      bool
	broker    Broker
}

func (m *Message) MetaData() map[string]interface{} {
	return make(map[string]interface{})
}

func (m *Message) Data() interface{} {
	return m.data
}

func (m *Message) QueueName() string {
	return m.queueName
}
