package workers

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
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
	QueueID     string `json:"QueueID" bson:"QueueID"`
	ID          string `json:"ID" bson:"ID"`
	queue       string
	Server      string `json:"Server" bson:"Server"`
	Handler     string `json:"Handler" bson:"Handler"`
	handler     JobFunc
	Concurrency int `json:"Concurrency" bson:"Concurrency"`
	concurrency int
	runners     []*taskRunner
	runnersLock sync.Mutex
	throttle    int
	Status      string    `json:"Status" bson:"Status"`
	Tag         []string  `json:"Tag" bson:"Tag"`
	StartedAt   int64     `json:"StartedAt" bson:"StartedAt"`
	PausedAt    time.Time `json:"PausedAt" bson:"PausedAt"`
	ResumedAt   time.Time `json:"ResumedAt" bson:"ResumedAt"`
	stop        chan bool
	StoppedAt   time.Time `json:"StoppedAt" bson:"StoppedAt"`
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

func newQueueWorker(queue string, wrk Worker, handler JobFunc, tag ...string) *Worker {
	if wrk.Concurrency <= 0 {
		wrk.Concurrency = 1
	}
	wrk.handler = handler
	wrk.queue = queue
	wrk.stop = make(chan bool)
	wrk.concurrency = wrk.Concurrency

	return &wrk
}

func (w *Worker) start(fetcher Fetcher) {
	w.runnersLock.Lock()
	if w.running {
		w.runnersLock.Unlock()
		return
	}
	w.running = true
	w.Status = STARTED
	w.StartedAt = time.Now().Unix()
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

func (w *Worker) Create() error {
	var err error
	if MG == nil {
		Connect(mongoURI, dbName)
	}
	var queue Queue
	query := bson.D{{Key: "ID", Value: w.QueueID}}
	record := MG.Db.Collection("queues").FindOne(context.Background(), query)
	err = record.Decode(&queue)
	if err != nil {
		panic(err)
		return err
	}
	w.Status = RUNNING
	w.StartedAt = time.Now().Unix()
	wrk := Worker{
		QueueID:     w.QueueID,
		ID:          w.ID,
		Server:      w.Server,
		handler:     DoWork,
		Concurrency: w.Concurrency,
		Status:      RUNNING,
		StartedAt:   time.Now().Unix(),
	}
	MG.Db.Collection("workers").InsertOne(context.Background(), w)
	if ManPool.Managers[w.QueueID] == nil {
		err = errors.New("workers for provided queue doesn't exists")
		return err
	}
	ManPool.Managers[w.QueueID].AddQueueWorker(queue.Name, wrk, DoWork)
	ManPool.Managers[w.QueueID].AdjustWorker(queue.Name, wrk)

	return nil
}

type Function struct {
	Map map[string]JobFunc
}

var FuncMapper = Function{}

func init() {
	FuncMapper.Map = map[string]JobFunc{}
}

func DoWork(message *Msg) error {
	mp, _ := message.Map()
	function := mp["class"].(string)
	err := FuncMapper.Map[function](message)
	if err != nil {
		panic(err)
	}
	return nil
}
