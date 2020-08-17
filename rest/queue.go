package rest

import (
	"github.com/gofiber/fiber"
	"github.com/segmentio/ksuid"
	"queue/workers"
	"time"
)

type QueueRequest struct {
	Addr        string
	WorkerID    string
	Concurrency int
}

func QueueHandlers(app fiber.Router) {
	queues := app.Group("/queues")
	queues.Post("/", AddQueue)
	queues.Post("/:id/worker/add", AddWorkerToQueue)
	queues.Post("/:id/worker/remove", RemoveWorkerFromQueue)
	queues.Post("/:id/worker/tune", TuneWorkerOnQueue)
	queues.Post("/:id/worker/tune", TuneWorkerOnQueue)
}

func AddQueue(c *fiber.Ctx) {
	var q workers.Queue
	c.BodyParser(&q)
	q.ID = ksuid.New().String()
	if q.MinWorkers == 0 {
		q.MinWorkers = 1
	}
	if q.MaxWorkers == 0 {
		q.MaxWorkers = 1
	}
	if q.MinPool == 0 {
		q.MinPool = 1
	}
	if q.MaxPool == 0 {
		q.MaxPool = 1
	}
	q.Status = workers.NOT_STARTED
	if q.StartImmediately {
		q.Status = workers.RUNNING
		q.StartedAt = time.Now().Unix()

	}
	q.Create(q.StartImmediately, DoWork)
	c.JSON(q)
}

func AddWorkerToQueue(c *fiber.Ctx) {
	queueId := c.Params("id")
	var w workers.Worker
	c.BodyParser(&w)
	w.ID = ksuid.New().String()
	w.QueueID = queueId
	w.Status = workers.NOT_STARTED
	if w.Concurrency == 0 {
		w.Concurrency = 1
	}
	w.Create()
	c.JSON(w)
}

func RemoveWorkerFromQueue(c *fiber.Ctx) {
	q := &workers.Queue{ID: c.Params("id")}
	var qr QueueRequest
	c.BodyParser(&qr)
	q.RemoveWorkersByAddress(qr.Addr)
	c.JSON(q)
}

func TuneWorkerOnQueue(c *fiber.Ctx) {
	q := &workers.Queue{ID: c.Params("id")}
	var qr QueueRequest
	c.BodyParser(&qr)
	q.Tune(qr.Concurrency)
}

func TuneWorkerOnQueueByWorkerId(c *fiber.Ctx) {
	q := &workers.Queue{ID: c.Params("id")}
	var qr QueueRequest
	c.BodyParser(&qr)
	q.TuneByWorkerId(qr.WorkerID, qr.Concurrency)
}
