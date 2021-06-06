package lika_queue

import (
	"github.com/lika_queue/brokers"
	"sync"
)

type QueueWorker struct {
	Broker   brokers.BrokerInterface
	Callable func(message brokers.MessageInterface)
	Worker   *Worker
}

func (qw *QueueWorker) Run() {
	qw.Worker.Callable = qw.consume
	qw.Worker.Run()
}

func (qw *QueueWorker) consume(queueName string, params map[string]interface{}) error {
	for {
		m, err := qw.Broker.Consume(queueName, params)

		if err != nil {
			return err
		}

		if m == nil && !qw.Broker.IsInfinite() {
			break
		}

		qw.Callable(m)
	}

	return nil
}

type Worker struct {
	Threads     int
	SyncThreads bool
	Callable    interface{}
}

func (w *Worker) Run() {
	if w.Threads > 1 {
		w.runMultiple()
	} else {
		w.Callable.(func())()
	}
}

func (w *Worker) runMultiple() {
	i := 0
	if w.SyncThreads {
		group := &sync.WaitGroup{}
		group.Add(w.Threads)

		for i < w.Threads {
			go w.runGroupThread(group)
			i++
		}

		group.Wait()
	} else {
		for i < w.Threads {
			go w.Callable.(func())()
			i++
		}
	}
}

func (w *Worker) runGroupThread(g *sync.WaitGroup) {
	defer g.Done()

	w.Callable.(func())()
}
