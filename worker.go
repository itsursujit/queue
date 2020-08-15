package main

import (
	"fmt"
	"github.com/segmentio/ksuid"
	"queue/workers"
)

func myJob(message *workers.Msg) error {
	fmt.Println(message.Jid())
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

func main() {
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

	if err != nil {
		fmt.Println(err)
	}

	// create a middleware chain with the default middlewares, and append myMiddleware
	mids := workers.DefaultMiddlewares().Append(myMiddleware)

	// pull messages from "myqueue" with concurrency of 10
	// this worker will not run myMiddleware, but will run the default middlewares
	manager.AddWorker("myqueue", 10, myJob)

	// pull messages from "myqueue2" with concurrency of 20
	// this worker will run the default middlewares and myMiddleware
	manager.AddWorker("myqueue2", 20, myJob, mids...)

	// pull messages from "myqueue3" with concurrency of 20
	// this worker will only run myMiddleware
	manager.AddWorker("myqueue3", 20, myJob, myMiddleware)

	// Blocks until process is told to exit via unix signal
	manager.Run()
	// stats will be available at http://localhost:8080/stats
	// workers.StartAPIServer(8080)
}
