package db

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	app *MongoDbOptions

	mDbClient *mongo.Client

	col *mongo.Collection
)

func ConnectToMongoDb(ops *MongoDbOptions) error {
	client, err := mongo.Connect(NewTimoutContext(10), options.Client().ApplyURI(ops.Uri))
	if err != nil {
		return err
	}
	err = client.Ping(NewTimoutContext(2), readpref.Primary())
	if err != nil {
		return err
	}
	app = ops
	mDbClient = client
	col = GetColRequests()
	return nil
}

func GetColRequests() *mongo.Collection {
	return mDbClient.Database(app.DbName).Collection(app.GetColName(ColRequests))
}