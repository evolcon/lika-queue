package brokers

type MessageInterface interface {
	GetData() interface{}
	GetMetaData() map[string]interface{}
	GetQueueName() string
}

func NewMessage(broker BrokerInterface, data interface{}, queueName string, metaData map[string]interface{}) MessageInterface {
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
	broker    BrokerInterface
}

func (m *Message) GetMetaData() map[string]interface{} {
	return make(map[string]interface{})
}

func (m *Message) GetData() interface{} {
	return m.data
}

func (m *Message) GetQueueName() string {
	return m.queueName
}
