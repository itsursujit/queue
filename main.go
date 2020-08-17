package main

import (
	"github.com/gofiber/fiber"
	"queue/rest"
	"queue/workers"
	"sync"
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

func myMiddleware(queue string, mgr *workers.Manager, next workers.JobFunc) workers.JobFunc {
	return func(message *workers.Msg) (err error) {
		// do something before each message is processed
		err = next(message)
		// do something after each message is processed
		return
	}
}

var managerPool *ManagerPool

func main() {
	app := fiber.New(&fiber.Settings{
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
	})
	rest.QueueHandlers(app)
	if err := app.Listen(":8080"); err != nil {
		println(err)
	}
}
