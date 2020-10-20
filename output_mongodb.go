package logbus

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoOutput struct {
	client      *mongo.Client
	dbName      string
	transformer TransformerFunc
}

func NewMongoOutput(host string, port int, username string, password string, dbName string,
	transformer TransformerFunc) (Output, error) {

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d", username, password, host, port)

	ctx1, cancel1 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel1()
	client, err := mongo.Connect(ctx1, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()
	if err := client.Ping(ctx2, readpref.Primary()); err != nil {
		return nil, err
	}

	o := MongoOutput{
		client:      client,
		dbName:      dbName,
		transformer: transformer,
	}
	return &o, nil
}

func (out *MongoOutput) Write(l *StdLog) error {
	lb, ok := out.transformer(l)
	if !ok {
		return nil
	}
	_, err := out.getC(lb.AppName).InsertOne(context.TODO(), l)
	return err
}

func (out *MongoOutput) Close() error {
	return out.client.Disconnect(context.TODO())
}

func (out *MongoOutput) getC(appName string) *mongo.Collection {
	return out.client.Database(out.dbName).Collection(fmt.Sprintf("log_%s", appName))
}
