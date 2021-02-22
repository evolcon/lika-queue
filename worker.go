package lika_queue

import "sync"

type QueueWorker struct {
	Queue    *Queue
	Callable func(message *Message)
	Worker   *Worker
}

func (qw *QueueWorker) Run() {
	qw.Worker.Callable = qw.consume
	qw.Worker.Run()
}

func (qw *QueueWorker) consume() {
	for {
		m := qw.Queue.Consume()

		if m == nil {
			break
		}

		qw.Callable(m)
	}
}

type Worker struct {
	Threads     int
	SyncThreads bool
	Callable    func()
}

func (w *Worker) Run() {
	if w.Threads > 1 {
		w.runMultiple()
	} else {
		w.Callable()
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

	w.Callable()
}
