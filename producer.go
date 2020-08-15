package main

import (
	"fmt"
	"github.com/segmentio/ksuid"
	"queue/workers"
)

func main() {
	// Create a manager, which manages workers
	producer, err := workers.NewProducer(workers.Options{
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
	for i := 0; i <= 100000; i++ {
		// Add a job to a queue
		producer.Enqueue("myqueue3", "Add", []int{1, 2})
	}

}
