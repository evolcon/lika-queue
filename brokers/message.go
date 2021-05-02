package brokers

type Message struct {
	Data      interface{}
	Retries   int
	QueueName string
	done      bool
	MQ        BrokerInterface
	id        int
}

func (m *Message) Close() {
	m.done = true
}

func (m *Message) Retrieve() {
	if m.done == false {
		m.Retries++
		m.MQ.PublishMessage(m.QueueName, m, nil)
	}
}
