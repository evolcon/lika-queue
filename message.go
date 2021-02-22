package lika_queue

type Message struct {
	Data    interface{}
	Retries int
	done    bool
	queue   *Queue
	id      int
}

func (m *Message) Close() {
	m.done = true
}

func (m *Message) Retrieve() {
	if m.done == false {
		m.Retries++
		m.queue.PublishMessage(m)
	}
}
