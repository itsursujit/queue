package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber"
	"github.com/segmentio/ksuid"
	"go.mongodb.org/mongo-driver/bson"
	"queue/workers"
	"sync"
	"time"
)

type ManagerPool struct {
	Managers map[string]*workers.Manager
	m        sync.Mutex
}
type QueueWorker struct {
	Queue    string
	Tag      string
	Worker   workers.Worker
	PoolSize uint
}

func SendEmail(message *workers.Msg) error {
	fmt.Println("I'm sending Email")
	return nil
}

func SendSMS(message *workers.Msg) error {
	fmt.Println("I'm sms sending")
	return nil
}

func myMiddleware(queue string, mgr *workers.Manager, next workers.JobFunc) workers.JobFunc {
	return func(message *workers.Msg) (err error) {
		// do something before each message is processed
		err = next(message)
		// do something after each message is processed
		return
	}
}

var managerPool *ManagerPool

const dbName = "fiber_test"
const mongoURI = "mongodb://localhost:27017/" + dbName

func main() {
	app := fiber.New(&fiber.Settings{
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
	})
	// functionList := []string{"SendEmail", "SendSMS"}
	userIds := []int{1, 2, 3, 4}
	queue := "myqueue"
	workers.Connect(mongoURI, dbName)
	for _, id := range userIds {
		queueId := ksuid.New().String()
		CreateWorkersForQueue(queue, queueId, id)
		StartWorkersForQueue(queueId)
		// TestProducer(queueId)
	}
	if err := app.Listen(":8080"); err != nil {
		println(err)
	}
}

func CreateWorkersForQueue(queue string, queueId string, userId int) {
	que := workers.Queue{
		Name:        fmt.Sprintf("%s:%v", queue, userId),
		Tag:         fmt.Sprintf("%v", userId),
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
		ID:          ksuid.New().String(),
		Tag:         []string{fmt.Sprintf("%v", userId)},
	}
	workers.MG.Db.Collection("workers").InsertOne(context.Background(), wrk)
	wrk = workers.Worker{
		QueueID:     queueId,
		Server:      "127.0.1.1",
		Handler:     "SendSMS",
		Status:      workers.NOT_STARTED,
		Concurrency: 100,
		ID:          ksuid.New().String(),
		Tag:         []string{fmt.Sprintf("%v", userId)},
	}
	workers.MG.Db.Collection("workers").InsertOne(context.Background(), wrk)
}

func StartWorkersForQueue(id string) error {
	managers := make(map[string]*workers.Manager)
	if managerPool == nil {
		managerPool = &ManagerPool{
			Managers: managers,
		}
	}
	// https://github.com/tsuru/monsterqueue
	// Create a manager, which manages workers
	manager, err := workers.NewManager(workers.Options{
		// location of redis instance
		ServerAddr: "localhost:6379",
		// instance of the database
		Database: 0,
		// number of connections to keep open with redis
		PoolSize: 30,
		// unique process id for this instance of workers (for proper recovery of inprogress jobs on crash)
		ProcessID: ksuid.New().String(),
	})
	var queue workers.Queue
	var wrkrs []workers.Worker
	query := bson.D{{Key: "ID", Value: id}}
	record := workers.MG.Db.Collection("queues").FindOne(context.Background(), query)
	err = record.Decode(&queue)
	if err != nil {
		return err
	}
	query = bson.D{{Key: "QueueID", Value: queue.ID}}
	cursor, err := workers.MG.Db.Collection("workers").Find(context.Background(), query)
	if err != nil {
		return err
	}
	err = cursor.All(context.Background(), &wrkrs)
	if err != nil {
		panic(err)
	}
	for _, wrk := range wrkrs {
		switch wrk.Handler {
		case "SendEmail":
			manager.AddWorker(queue.Name, wrk.Concurrency, SendEmail)
		case "SendSMS":
			manager.AddWorker(queue.Name, wrk.Concurrency, SendSMS)
		}
	}
	managerPool.m.Lock()
	managerPool.Managers[id] = manager
	managerPool.m.Unlock()
	go managerPool.Managers[id].Run()
	return nil
}
