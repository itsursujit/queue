package main

import (
	"context"
	"fmt"
	"github.com/segmentio/ksuid"
	"queue/workers"
	"sync"
	"time"
)

type ManagerPool struct {
	Managers map[string]workers.Manager
	m        sync.Mutex
}
type QueueWorker struct {
	Queue    string
	Tag      string
	Worker   workers.Worker
	PoolSize uint
}

func SendEmail(args interface{}) {

}
func SendSMS(args interface{}) {

}

var managerPool *ManagerPool

const dbName = "fiber_test"
const mongoURI = "mongodb://localhost:27017/" + dbName

func main() {
	// functionList := []string{"SendEmail", "SendSMS"}
	userIds := []int{1, 2, 3, 4}
	queue := "myqueue"
	workers.Connect(mongoURI, dbName)
	for _, id := range userIds {
		queueId := ksuid.New().String()
		que := workers.Queue{
			Name:        queue,
			Tag:         fmt.Sprintf("%v", id),
			Status:      workers.NOT_STARTED,
			ID:          queueId,
			IsDedicated: true,
			StartedAt:   time.Now().Unix(),
		}
		workers.MG.Db.Collection("queues").InsertOne(context.Background(), que)
		wrk := workers.Worker{
			QueueID:     queueId,
			Server:      "127.0.0.1",
			Handler:     "SendEmail",
			Status:      workers.NOT_STARTED,
			Concurrency: 100,
			Tag:         []string{fmt.Sprintf("%v", id)},
		}
		workers.MG.Db.Collection("workers").InsertOne(context.Background(), wrk)
		wrk = workers.Worker{
			QueueID:     queueId,
			Server:      "127.0.1.1",
			Handler:     "SendEmail",
			Status:      workers.NOT_STARTED,
			Concurrency: 100,
			Tag:         []string{fmt.Sprintf("%v", id)},
		}
		workers.MG.Db.Collection("workers").InsertOne(context.Background(), wrk)
	}

}

func CreateWorkersForQueue(id string) {

}
