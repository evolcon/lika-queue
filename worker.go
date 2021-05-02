package lika_queue

import (
	"github.com/lika_queue/brokers"
	"sync"
)

type QueueWorker struct {
	Queue    brokers.BrokerInterface
	Callable func(message *brokers.Message)
	Worker   *Worker
}

func (qw *QueueWorker) Run() {
	qw.Worker.Callable = qw.consume
	qw.Worker.Run()
}

func (qw *QueueWorker) consume(queueName string, params map[string]interface{}) {
	for {
		m := qw.Queue.Consume(queueName, params)

		if m == nil && !qw.Queue.IsInfinite() {
			break
		}

		qw.Callable(m)
	}
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
