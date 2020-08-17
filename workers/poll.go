package workers

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

func PollTasks() {
	var tasks []Task
	q := bson.D{{}}
	cursor, err := MG.Db.Collection("jobs").Find(context.Background(), q)
	if err != nil {
		panic(err)
	}
	if err := cursor.All(context.Background(), &tasks); err != nil {
		panic(err)
	}
}
