package workers

import (
	"context"
	"fmt"
	"github.com/segmentio/ksuid"
	"go.mongodb.org/mongo-driver/bson"
)

type Queue struct {
	Name                string     `json:"Name" bson:"Name"`
	Tag                 string     `json:"Tag" bson:"Tag"`
	Workers             []Worker   `json:"Workers" bson:"Workers"`
	Status              string     `json:"Status" bson:"Status"`
	ID                  string     `json:"ID" bson:"ID"`
	StartImmediately    bool       `json:"StartImmediately" bson:"StartImmediately"`
	QueueStats          QueueStats `json:"QueueStats" bson:"QueueStats"`
	StartedAt           int64      `json:"StartedAt" bson:"StartedAt"`
	MinWorkers          int        `json:"MinWorkers" bson:"MinWorkers"`
	MaxWorkers          int        `json:"MaxWorkers" bson:"MaxWorkers"`
	MinPool             int        `json:"MinPool" bson:"MinPool"`
	MaxPool             int        `json:"MaxPool" bson:"MaxPool"`
	TriggerInitWorkers  bool       `json:"TriggerInitWorkers" bson:"TriggerInitWorkers"`
	TriggerCloseWorkers bool       `json:"TriggerCloseWorkers" bson:"TriggerCloseWorkers"`
	ParentID            string     `json:"ParentID" bson:"ParentID"`
	jobFunc             JobFunc
}

type JobStats struct {
	NotStarted uint64 `json:"NotStarted" bson:"NotStarted"`
	Started    uint64 `json:"Started" bson:"Started"`
	InProgress uint64 `json:"InProgress" bson:"InProgress"`
	Paused     uint64 `json:"Paused" bson:"Paused"`
	Completed  uint64 `json:"Completed" bson:"Completed"`
	Failed     uint64 `json:"Failed" bson:"Failed"`
	Expired    uint64 `json:"Expired" bson:"Expired"`
	Canceled   uint64 `json:"Canceled" bson:"Canceled"`
}

type WorkerStats struct {
	NotStarted uint64 `json:"NotStarted" bson:"NotStarted"`
	Running    uint64 `json:"Running" bson:"Running"`
	Paused     uint64 `json:"Paused" bson:"Paused"`
	Stopped    uint64 `json:"Stopped" bson:"Stopped"`
	Killed     uint64 `json:"Killed" bson:"Killed"`
}

type QueueStats struct {
	JobStats    JobStats    `json:"JobStats" bson:"JobStats"`
	WorkerStats WorkerStats `json:"WorkerStats" bson:"WorkerStats"`
}

func (q *Queue) Create(startImmediate bool, doWork JobFunc) {
	if MG == nil {
		Connect(mongoURI, dbName)
	}
	q.jobFunc = doWork
	MG.Db.Collection("queues").InsertOne(context.Background(), q)
	if startImmediate {
		wrk := Worker{
			QueueID:     q.ID,
			Server:      "127.0.0.1",
			Handler:     "SendEmail",
			Status:      RUNNING,
			Concurrency: q.MinPool,
			ID:          ksuid.New().String(),
			Tag:         []string{q.Tag},
		}
		MG.Db.Collection("workers").InsertOne(context.Background(), wrk)
		q.StartWorkers()
	}
}

func (q *Queue) StartWorkers() error {
	if MG == nil {
		Connect(mongoURI, dbName)
	}
	managers := make(map[string]*Manager)
	if ManPool == nil {
		ManPool = &ManagerPool{
			Managers: managers,
		}
	}
	// https://github.com/tsuru/monsterqueue
	// Create a manager, which manages workers
	manager, err := NewManager(Options{
		// location of redis instance
		ServerAddr: "localhost:6379",
		// instance of the database
		Database: 0,
		// number of connections to keep open with redis
		PoolSize: 30,
		// unique process id for this instance of workers (for proper recovery of inprogress jobs on crash)
		ProcessID: ksuid.New().String(),
	})
	var queue Queue
	var wrkrs []Worker
	query := bson.D{{Key: "ID", Value: q.ID}}
	record := MG.Db.Collection("queues").FindOne(context.Background(), query)
	err = record.Decode(&queue)
	if err != nil {
		return err
	}
	query = bson.D{{Key: "QueueID", Value: queue.ID}}
	cursor, err := MG.Db.Collection("workers").Find(context.Background(), query)
	if err != nil {
		return err
	}
	err = cursor.All(context.Background(), &wrkrs)
	if err != nil {
		panic(err)
	}
	for _, wrk := range wrkrs {
		manager.AddQueueWorker(queue.Name, wrk, q.jobFunc)
	}
	ManPool.m.Lock()
	ManPool.Managers[q.ID] = manager
	ManPool.m.Unlock()
	go ManPool.Managers[q.ID].Run()
	fmt.Println("Running")
	return nil
}

func (q *Queue) RemoveWorkersByAddress(addr ...string) {
	wrkSrv := ""
	if len(addr) > 0 {
		wrkSrv = addr[0]
	}
	if ManPool != nil {
		if ManPool.Managers[q.ID].IsRunning() {
			for _, w := range ManPool.Managers[q.ID].GetWorkers() {
				if w.Server == wrkSrv {
					ManPool.Managers[q.ID].StopWorker(w.ID)
				}
			}
		}
	}
}

func (q *Queue) Tune(concurrency int) error {
	ManPool.Managers[q.ID].Tune(concurrency)
	return nil
}

func (q *Queue) TuneByWorkerId(id string, concurrency int) error {
	ManPool.Managers[q.ID].TuneWorker(id, concurrency)
	return nil
}
