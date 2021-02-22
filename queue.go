package lika_queue

import (
	"sync"
	"time"
)

func New() *Queue {
	return &Queue{
		messages: make(map[int]*Message),
	}
}

func NewInfinite(duration time.Duration) *Queue {
	return &Queue{
		messages: make(map[int]*Message),
		Infinite: true,
		Duration: duration,
	}
}

type Queue struct {
	Infinite bool
	Duration time.Duration // duration time in milliseconds to waiting for a new message
	messages map[int]*Message
	offset   int
	lock     sync.Mutex
}

func (q *Queue) Publish(message interface{}) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.offset++
	q.messages[q.offset] = &Message{Data: message, id: q.offset, queue: q}
}

func (q *Queue) PublishMessage(message *Message) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.messages[message.id] = message
}

func (q *Queue) Consume() *Message {
	for {
		q.lock.Lock()

		if !q.HasMessages() {
			q.lock.Unlock()

			if q.Infinite {
				time.Sleep(q.Duration * time.Millisecond)
				continue
			}

			return nil
		}

		defer q.lock.Unlock()

		return q.RemoveMessage(q.CurrentIndex())
	}
}

func (q *Queue) RemoveMessage(key int) *Message {
	defer delete(q.messages, key)

	return q.messages[key]
}

func (q *Queue) HasMessages() bool {
	return q.Len() > 0
}

func (q *Queue) Len() int {
	return len(q.messages)
}

func (q *Queue) CurrentIndex() int {
	return q.offset - len(q.messages) + 1
}
