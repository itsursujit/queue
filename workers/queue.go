package workers

type Queue struct {
	Name       string   `json:"Name"`
	Tag        string   `json:"Tag"`
	Workers    []Worker `json:"Workers"`
	Status     string   `json:"Status"`
	ID         string   `json:"ID"`
	QueueStats QueueStats
	StartedAt  int64
}

type JobStats struct {
	NotStarted uint64 `json:"NotStarted"`
	Started    uint64 `json:"Started"`
	InProgress uint64 `json:"InProgress"`
	Paused     uint64 `json:"Paused"`
	Completed  uint64 `json:"Completed"`
	Failed     uint64 `json:"Failed"`
	Expired    uint64 `json:"Expired"`
	Canceled   uint64 `json:"Canceled"`
}

type WorkerStats struct {
	NotStarted uint64 `json:"NotStarted"`
	Running    uint64 `json:"Running"`
	Paused     uint64 `json:"Paused"`
	Stopped    uint64 `json:"Stopped"`
	Killed     uint64 `json:"Killed"`
}

type QueueStats struct {
	JobStats    JobStats    `json:"JobStats"`
	WorkerStats WorkerStats `json:"WorkerStats"`
}
