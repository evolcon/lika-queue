package brokers

type MessageInterface interface {
	GetData() interface{}
	GetMetaData() map[string]interface{}
	GetRetries() int
	GetQueueName() string
	Close()
	Retrieve() error
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

func (m *Message) GetRetries() int {
	return m.retries
}

func (m *Message) GetQueueName() string {
	return m.queueName
}

func (m *Message) Close() {
	m.done = true
}

func (m *Message) Retrieve() error {
	if m.done == false {
		m.retries++
		return m.broker.PublishMessage(m.queueName, m, nil)
	}

	return nil
}
