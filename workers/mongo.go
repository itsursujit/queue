package workers

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoInstance contains the Mongo client and database objects
type MongoInstance struct {
	Client         *mongo.Client
	Db             *mongo.Database
	CollectionName string
}

var MG *MongoInstance
var _ Persist = &MongoInstance{}

type Persist interface {

	// General queue operations
	CreateQueue(queue string) error
	ListMessages(queue string) ([]EnqueueData, error)
	EnqueueMessage(queue string, priority float64, message EnqueueData) error
	EnqueueMessageNow(queue string, message EnqueueData) error
	// Special purpose queue operations
	EnqueueScheduledMessage(priority float64, message EnqueueData) error
}

// Connect configures the MongoDB client and initializes the database connection.
// Source: https://www.mongodb.com/blog/post/quick-start-golang--mongodb--starting-and-setup
func Connect(mongoURI string, dbName string) error {
	if MG != nil {
		return nil
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	db := client.Database(dbName)

	if err != nil {
		return err
	}

	MG = &MongoInstance{
		Client:         client,
		Db:             db,
		CollectionName: "jobs",
	}

	return nil
}

func (r *MongoInstance) EnqueueMessage(queue string, priority float64, message EnqueueData) error {
	collection := r.Db.Collection(r.CollectionName)
	_, err := collection.InsertOne(context.TODO(), message)
	if err != nil {
		return err
	}
	return nil
}

func (r *MongoInstance) EnqueueScheduledMessage(priority float64, message EnqueueData) error {
	collection := r.Db.Collection(r.CollectionName)
	_, err := collection.InsertOne(context.TODO(), message)
	if err != nil {
		return err
	}
	return nil
}

func (r *MongoInstance) EnqueueMessageNow(queue string, message EnqueueData) error {
	collection := r.Db.Collection(r.CollectionName)
	_, err := collection.InsertOne(context.TODO(), message)
	if err != nil {
		return err
	}
	return nil
}

func (r *MongoInstance) CreateQueue(queue string) error {
	return nil
}

func (r *MongoInstance) ListMessages(queue string) ([]EnqueueData, error) {
	return nil, nil
}
