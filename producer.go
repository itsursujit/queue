package main

import (
	"fmt"
	"github.com/segmentio/ksuid"
	"queue/workers"
)

const dbName = "fiber_test"
const mongoURI = "mongodb://localhost:27017/" + dbName

func main() {
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
		fmt.Println(err)
	}
	for i := 0; i <= 10; i++ {
		// Add a job to a queue
		producer.Enqueue("myqueue:3", "Add", []int{1, 2})
	}

}
