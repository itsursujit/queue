package rest

import "github.com/gofiber/fiber"

func TaskHandlers(app fiber.Router) {
	tasks := app.Group("/tasks")
	tasks.Get("/", ListTask)
	tasks.Get("/:queue", TaskListByQueueName)
	tasks.Get("/:tag", TaskListByQueueName)
}

func CreateTask(c *fiber.Ctx) {

}

func ViewTask(c *fiber.Ctx) {

}

func DeleteTask(c *fiber.Ctx) {

}

func ListTask(c *fiber.Ctx) {

}

func TaskListByQueueName(c *fiber.Ctx) {

}

func ListTaskByTag(c *fiber.Ctx) {

}

func TaskStats(c *fiber.Ctx) {

}

func TaskStatsByQueueName(c *fiber.Ctx) {

}

func TaskStatsByTag(c *fiber.Ctx) {

}
