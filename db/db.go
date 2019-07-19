package db

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// New creates a new mongodb connection and tests it
func New(url string, db string) (*mongo.Database, error) {
	// register custom codec registry to handle empty interfaces
	rb := bson.NewRegistryBuilder()
	rb.RegisterTypeMapEntry(bsontype.EmbeddedDocument, reflect.TypeOf(bson.M{}))

	// create mongodb connection
	client, err := mongo.NewClient(options.Client().ApplyURI(url).SetRegistry(rb.Build()))
	if err != nil {
		return nil, fmt.Errorf("db new client error: %s", err)
	}

	// connect to the mongodb
	ctxc, donec := context.WithTimeout(context.Background(), 10*time.Second)
	defer donec()
	if err := client.Connect(ctxc); err != nil {
		return nil, fmt.Errorf("db connection error: %s", err)
	}

	// is the mongo really there?
	ctxp, donep := context.WithTimeout(context.Background(), 2*time.Second)
	defer donep()
	if err := client.Ping(ctxp, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("db ping error: %s", err)
	}

	return client.Database(db), nil
}
