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
}

func AddQueue(c *fiber.Ctx) {
	var q workers.Queue
	err := c.BodyParser(&q)
	if err != nil {
		panic(err)
	}
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
	q.Create(q.StartImmediately, workers.DoWork)
	err = c.JSON(q)
	if err != nil {
		panic(err)
	}
}

func AddWorkerToQueue(c *fiber.Ctx) {
	queueId := c.Params("id")
	var w workers.Worker
	err := c.BodyParser(&w)
	if err != nil {
		panic(err)
	}
	w.ID = ksuid.New().String()
	w.QueueID = queueId
	w.Status = workers.NOT_STARTED
	if w.Concurrency == 0 {
		w.Concurrency = 1
	}
	err = w.Create()
	if err != nil {
		panic(err)
	}
	err = c.JSON(w)
	if err != nil {
		panic(err)
	}
}

func RemoveWorkerFromQueue(c *fiber.Ctx) {
	q := &workers.Queue{ID: c.Params("id")}
	var qr QueueRequest
	err := c.BodyParser(&qr)
	if err != nil {
		panic(err)
	}
	q.RemoveWorkersByAddress(qr.Addr)
	err = c.JSON(q)
}

func TuneWorkerOnQueue(c *fiber.Ctx) {
	q := &workers.Queue{ID: c.Params("id")}
	var qr QueueRequest
	err := c.BodyParser(&qr)
	if err != nil {
		panic(err)
	}
	if qr.WorkerID != "" {
		err = q.TuneByWorkerId(qr.WorkerID, qr.Concurrency)
		if err != nil {
			panic(err)
		}
	} else {
		err = q.Tune(qr.Concurrency)
		if err != nil {
			panic(err)
		}
	}
}

func TuneWorkerOnQueueByWorkerId(c *fiber.Ctx) {
	q := &workers.Queue{ID: c.Params("id")}
	var qr QueueRequest
	err := c.BodyParser(&qr)
	if err != nil {
		panic(err)
	}
	err = q.TuneByWorkerId(qr.WorkerID, qr.Concurrency)
	if err != nil {
		panic(err)
	}
}
