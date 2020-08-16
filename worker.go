package main

import (
	"fmt"
	"github.com/gofiber/fiber"
	"github.com/segmentio/ksuid"
	"queue/workers"
)

func myJob(message *workers.Msg) error {
	fmt.Println(message.Jid())
	return nil
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
	manager.AddWorker("myqueue:3", 20, myJob)
	// Blocks until process is told to exit via unix signal

	app := fiber.New(&fiber.Settings{
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
	})
	go manager.Run()
	if err := app.Listen(":8081"); err != nil {
		println(err)
	}
	// stats will be available at http://localhost:8080/stats
	// workers.StartAPIServer(8080)
}
