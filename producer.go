package main

import (
	"fmt"
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
		ProcessID: "1",
	})

	if err != nil {
		fmt.Println(err)
	}
	go func() {
		for i := 0; i <= 100000; i++ {
			// Add a job to a queue
			producer.Enqueue("myqueue3", "Add", []int{1, 2})
		}
	}()

}
