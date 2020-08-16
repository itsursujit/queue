package main

import (
	"github.com/gofiber/fiber"
	"queue/rest"
	"queue/workers"
)

type ManagerPool struct {
}

var ManagerPool workers.Manager

func main() {
	app := fiber.New(&fiber.Settings{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
	})
	group := app.Group("/")
	rest.QueueHandlers(group)
	rest.TaskHandlers(group)
	rest.WorkerHandlers(group)

	if err := app.Listen(":8080"); err != nil {
		println(err)
	}
}
