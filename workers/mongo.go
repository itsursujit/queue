package workers

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
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

const dbName = "fiber_test"
const mongoURI = "mongodb://localhost:27017/" + dbName

// Connect configures the MongoDB client and initializes the database connection.
// Source: https://www.mongodb.com/blog/post/quick-start-golang--mongodb--starting-and-setup
func Connect(mongoURI string, dbName string) error {
	if MG != nil {
		return nil
	}
	client, err := mongo.NewClient(options.Client().SetMaxPoolSize(200).ApplyURI(mongoURI))

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
		Client: client,
		Db:     db,
	}

	return nil
}

func (r *MongoInstance) EnqueueMessage(queue string, priority float64, message EnqueueData) error {
	collection := r.Db.Collection(r.CollectionName)
	_, err := collection.InsertOne(context.Background(), message)
	if err != nil {
		return err
	}
	return nil
}

func (r *MongoInstance) EnqueueScheduledMessage(priority float64, message EnqueueData) error {
	collection := r.Db.Collection(r.CollectionName)
	_, err := collection.InsertOne(context.Background(), message)
	if err != nil {
		return err
	}
	return nil
}

func (r *MongoInstance) EnqueueMessageNow(queue string, message EnqueueData) error {
	collection := r.Db.Collection(r.CollectionName)
	_, err := collection.InsertOne(context.Background(), message)
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

// ChangeStream defines what to watch? client, database or collection
type ChangeStream struct {
	Collection string
	Database   string
	Pipeline   []bson.D
}

type callback func(bson.M)

// SetCollection sets collection
func (cs *ChangeStream) SetCollection(collection string) {
	cs.Collection = collection
}

// SetDatabase sets database
func (cs *ChangeStream) SetDatabase(database string) {
	cs.Database = database
}

// SetPipeline sets pipeline
func (cs *ChangeStream) SetPipeline(pipeline []bson.D) {
	cs.Pipeline = pipeline
}

// NewChangeStream gets a new ChangeStream
func NewChangeStream() *ChangeStream {
	return &ChangeStream{}
}

// Watch prints oplogs in JSON format
func (cs *ChangeStream) Watch(client *mongo.Client, cb callback) {
	var err error
	var ctx = context.Background()
	var cur *mongo.ChangeStream
	fmt.Println("pipeline", cs.Pipeline)
	opts := options.ChangeStream()
	opts.SetFullDocument("updateLookup")
	if cs.Collection != "" && cs.Database != "" {
		fmt.Println("Watching", cs.Database+"."+cs.Collection)
		var coll = client.Database(cs.Database).Collection(cs.Collection)
		if cur, err = coll.Watch(ctx, cs.Pipeline, opts); err != nil {
			panic(err)
		}
	} else if cs.Database != "" {
		fmt.Println("Watching", cs.Database)
		var db = client.Database(cs.Database)
		if cur, err = db.Watch(ctx, cs.Pipeline, opts); err != nil {
			panic(err)
		}
	} else {
		fmt.Println("Watching all")
		if cur, err = client.Watch(ctx, cs.Pipeline, opts); err != nil {
			panic(err)
		}
	}

	defer cur.Close(ctx)
	var doc bson.M
	for cur.Next(ctx) {
		if err = cur.Decode(&doc); err != nil {
			log.Fatal(err)
		}
		cb(doc)
	}
	if err = cur.Err(); err != nil {
		log.Fatal(err)
	}
}
