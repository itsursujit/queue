package main

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"queue/workers"
)

func main() {
	const dbName = "fiber_test"
	const mongoURI = "mongodb://localhost:27017/" + dbName + "?replicaSet=replset"
	if err := workers.Connect(mongoURI, dbName); err != nil {
		log.Fatal(err)
	}
	cs := workers.ChangeStream{
		Collection: "jobs",
		Database:   dbName,
		Pipeline:   mongo.Pipeline{},
	}
	cs.Watch(workers.MG.Client, func(m bson.M) {
		fmt.Println("I found you")
		fmt.Println(m)
	})
}
