package workers

type Queue struct {
	Name             string   `json:"Name" bson:"Name"`
	Tag              string   `json:"Tag" bson:"Tag"`
	Workers          []Worker `json:"Workers" bson:"Workers"`
	Status           string   `json:"Status" bson:"Status"`
	ID               string   `json:"ID" bson:"ID"`
	IsDedicated      bool     `json:"IsDedicated" bson:"IsDedicated"`
	StartImmediately bool     `json:"StartImmediately" bson:"StartImmediately"`
	QueueStats       QueueStats
	StartedAt        int64
}

type JobStats struct {
	NotStarted uint64 `json:"NotStarted" bson:"NotStarted"`
	Started    uint64 `json:"Started" bson:"Started"`
	InProgress uint64 `json:"InProgress" bson:"InProgress"`
	Paused     uint64 `json:"Paused" bson:"Paused"`
	Completed  uint64 `json:"Completed" bson:"Completed"`
	Failed     uint64 `json:"Failed" bson:"Failed"`
	Expired    uint64 `json:"Expired" bson:"Expired"`
	Canceled   uint64 `json:"Canceled" bson:"Canceled"`
}

type WorkerStats struct {
	NotStarted uint64 `json:"NotStarted" bson:"NotStarted"`
	Running    uint64 `json:"Running" bson:"Running"`
	Paused     uint64 `json:"Paused" bson:"Paused"`
	Stopped    uint64 `json:"Stopped" bson:"Stopped"`
	Killed     uint64 `json:"Killed" bson:"Killed"`
}

type QueueStats struct {
	JobStats    JobStats    `json:"JobStats" bson:"JobStats"`
	WorkerStats WorkerStats `json:"WorkerStats" bson:"WorkerStats"`
}
