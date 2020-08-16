package workers

import (
	"sync"
	"time"
)

const (
	NOT_STARTED = "NOT_STARTED"
	STARTED     = "STARTED"
	RUNNING     = "RUNNING"
	IN_PROGRESS = "IN_PROGRESS"
	PROCESSING  = "PROCESSING"
	PAUSED      = "PAUSED"
	CANCELED    = "CANCELED"
	EXPIRED     = "EXPIRED"
	STOPPED     = "STOPPED"
	FAILED      = "FAILED"
	COMPLETED   = "COMPLETED"
)

type Worker struct {
	QueueID     string
	queue       string
	Server      string
	Handler     string
	handler     JobFunc
	Concurrency int
	concurrency int
	runners     []*taskRunner
	runnersLock sync.Mutex
	throttle    int
	Status      string
	Tag         []string
	StartedAt   time.Time
	PausedAt    time.Time
	ResumedAt   time.Time
	stop        chan bool
	StoppedAt   time.Time
	running     bool
}

func newWorker(queue string, concurrency int, handler JobFunc, tag ...string) *Worker {
	if concurrency <= 0 {
		concurrency = 1
	}
	w := &Worker{
		queue:       queue,
		handler:     handler,
		concurrency: concurrency,
		stop:        make(chan bool),
		Tag:         tag,
	}
	return w
}

func (w *Worker) start(fetcher Fetcher) {
	w.runnersLock.Lock()
	if w.running {
		w.runnersLock.Unlock()
		return
	}
	w.running = true
	w.Status = STARTED
	w.StartedAt = time.Now()
	defer func() {
		w.runnersLock.Lock()
		w.running = false
		w.runnersLock.Unlock()
	}()

	var wg sync.WaitGroup
	wg.Add(w.concurrency)

	go fetcher.Fetch()

	done := make(chan *Msg)
	w.runners = make([]*taskRunner, w.concurrency)
	for i := 0; i < w.concurrency; i++ {
		r := newTaskRunner(w.handler)
		w.runners[i] = r
		w.Status = RUNNING
		go func() {
			r.work(fetcher.Messages(), done, fetcher.Ready())
			wg.Done()
		}()
	}
	exit := make(chan bool)
	go func() {
		wg.Wait()
		close(exit)
	}()

	// Now that we're all set up, unlock so that stats can check.
	w.runnersLock.Unlock()

	for {
		select {
		case msg := <-done:
			if msg.ack {
				fetcher.Acknowledge(msg)
			}
		case <-w.stop:
			if !fetcher.Closed() {
				fetcher.Close()

				// we need to relock the runners so we can shut this down
				w.runnersLock.Lock()
				for _, r := range w.runners {
					r.quit()
				}
				w.runnersLock.Unlock()
			}
		case <-exit:
			return
		}
	}
}

func (w *Worker) quit() {
	w.runnersLock.Lock()
	defer w.runnersLock.Unlock()
	if w.running {
		w.Status = STOPPED
		w.StoppedAt = time.Now()
		w.stop <- true
	}
}

func (w *Worker) inProgressMessages() []*Msg {
	w.runnersLock.Lock()
	defer w.runnersLock.Unlock()
	var res []*Msg
	for _, r := range w.runners {
		if m := r.inProgressMessage(); m != nil {
			res = append(res, m)
		}
	}
	return res
}
