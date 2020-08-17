package main

import (
	"context"
	"fmt"
	"github.com/segmentio/ksuid"
	"go.mongodb.org/mongo-driver/bson"
	"queue/workers"
)

func main() {
	queueId := "1gEsL06XZBLUli8ji5XLdc0glsW"
	TestProducer(queueId)
}

func TestProducer(id string) error {
	const dbName = "fiber_test"
	const mongoURI = "mongodb://localhost:27017/" + dbName
	workers.Connect(mongoURI, dbName)
	var queue workers.Queue
	query := bson.D{{Key: "ID", Value: id}}
	record := workers.MG.Db.Collection("queues").FindOne(context.Background(), query)
	err := record.Decode(&queue)
	if err != nil {
		panic(err)
	}
	// Create a manager, which manages workers
	producer, err := workers.NewProducer(workers.Options{
		PersistentAddr: mongoURI,
		PersistentDB:   dbName,
		// location of redis instance
		ServerAddr: "localhost:6379",
		// instance of the database
		Database: 0,
		// number of connections to keep open with redis
		PoolSize: 30,
		// unique process id for this instance of workers (for proper recovery of inprogress jobs on crash)
		ProcessID: ksuid.New().String(),
	})
	if err != nil {
		panic(err)
	}
	for i := 0; i <= 10; i++ {
		fmt.Println(queue.Name)
		// _, err = producer.Enqueue(queue.Name, "SendEmail", []int{1, 2})
		_, err = producer.Enqueue("Test", "SendSMS", []int{1, 2})
	}
	return nil
}
