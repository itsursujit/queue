package rest

import (
	"github.com/gofiber/fiber"
	"queue/workers"
)

func QueueHandlers(app fiber.Router) {
	queues := app.Group("/queues")
	queues.Get("/", ListQueue)
	queues.Post("/", AddQueue)
}

func ListQueue(c *fiber.Ctx) {

}

func AddQueue(c *fiber.Ctx) {
	var q workers.Queue
	c.BodyParser(&q)
	c.JSON(q)
}

func DeleteQueue(c *fiber.Ctx) {

}

func AddWorkerToQueue(c *fiber.Ctx) {

}

func ListQueueByTag(c *fiber.Ctx) {

}

func ListQueueById(c *fiber.Ctx) {

}
