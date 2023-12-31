package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	//same name for both, but will keep them to make it right
	databaseName = "logs"
	collectionName = "logs"	
)

var client *mongo.Client

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID string `bson:"_id,omitempty" json:"id,omitempty"`
	Name string `bson:"name" json:"name"`
	Data string `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func New(mongo *mongo.Client) Models {
	client = mongo

	return Models{
		LogEntry: LogEntry{},
	}
}

func (l *LogEntry) Insert(e LogEntry) error {
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	//defering cancel in case things take longer than 15 seconds
	//defer cancel()

	collection := client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name: e.Name,
		Data: e.Data,
		CreatedAt: time.Now(),
	})

	if err != nil {
		log.Println("error inserting log:", err)
		return err
	}

	return nil
}

func (l *LogEntry) All() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	//defering cancel in case things take longer than 15 seconds
	defer cancel()

	coll := client.Database(databaseName).Collection(collectionName)

	opts := options.Find()
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	//Mongo cursor is the collection returned in the Find function
	cursor, err := coll.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("Error finding all logs", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*LogEntry

	for cursor.Next(ctx) {
		var l LogEntry
		err := cursor.Decode(&l)
		if err != nil {
			log.Println("Error decoding log into slice: ", err)
			return nil, err
		} else {
			logs = append(logs, &l)
		}		
	}

	return logs, nil
}

func (l *LogEntry) GetOne(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	//defering cancel in case things take longer than 15 seconds
	defer cancel()

	coll := client.Database(databaseName).Collection(collectionName)

	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var le LogEntry
	//need to decode the single result to get to the item
	err = coll.FindOne(ctx, bson.M{"_id": docId}).Decode(&le)
	if err != nil {
		log.Println("error on finding one:", err)
		return nil, err
	}

	return &le, nil
}

func (l *LogEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	//defering cancel in case things take longer than 15 seconds
	defer cancel()

	coll := client.Database(databaseName).Collection(collectionName)

	if err := coll.Drop(ctx); err != nil {
		return err
	}
	
	return nil	
}

// following through the course, I don't think log entries have to be updated at all
func (l *LogEntry) Update() (*mongo.UpdateResult, error) {	
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	//defering cancel in case things take longer than 15 seconds
	defer cancel()

	coll := client.Database(databaseName).Collection(collectionName)

	docId, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		return nil, err
	}

	res, err := coll.UpdateOne(
		ctx,
		bson.M{"_id": docId},
		bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "name", Value: l.Name},
				{Key: "data", Value: l.Data},
				{Key: "updated_at", Value: time.Now()},
			}},
		},
	)

	if err != nil {
		return nil, err
	}

	return res, nil
}