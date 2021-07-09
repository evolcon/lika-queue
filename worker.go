package lika_queue

import (
	"sync"
	"time"
)

type QueueWorker struct {
	QueueName       string
	Duration        time.Duration // milliseconds. default 100 ms
	ConsumingParams map[string]interface{}
	Broker          Broker
	Callable        func(message MessageData)
	Worker          *Worker
	IsInfinite      bool // if set false, message consuming will stop after getting first nil result from broker
}

func (qw *QueueWorker) Run() {
	if qw.IsInfinite && qw.Duration <= 0 {
		qw.Duration = 100
	}
	qw.Worker.Callable = qw.consume
	qw.Worker.Run()
}

func (qw *QueueWorker) consume() error {
	for {
		m, err := qw.Broker.Consume(qw.QueueName, qw.ConsumingParams)

		if err != nil {
			return err
		}

		if m == nil {
			if qw.IsInfinite {
				time.Sleep(qw.Duration * time.Millisecond)
				continue
			} else {
				break
			}
		}

		qw.Callable(m)
	}

	return nil
}

type Worker struct {
	Threads     int
	SyncThreads bool
	Callable    func() error
}

func (w *Worker) Run() {
	if w.Threads > 1 {
		w.runMultiple()
	} else {
		_ = w.Callable()
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
			go w.Callable()
			i++
		}
	}
}

func (w *Worker) runGroupThread(g *sync.WaitGroup) {
	defer g.Done()

	_ = w.Callable()
}
