package db

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/I1820/I1820/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const connectionTimeout = 10 * time.Second

// New creates a new mongodb connection and tests it.
func New(cfg config.Database) (*mongo.Database, error) {
	// register custom codec registry to handle empty interfaces
	rb := bson.NewRegistryBuilder()
	rb.RegisterTypeMapEntry(bsontype.EmbeddedDocument, reflect.TypeOf(bson.M{}))

	// create mongodb connection
	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.URL).SetRegistry(rb.Build()))
	if err != nil {
		return nil, fmt.Errorf("db new client error: %w", err)
	}

	{
		// connect to the mongodb
		ctx, done := context.WithTimeout(context.Background(), connectionTimeout)
		defer done()

		if err := client.Connect(ctx); err != nil {
			return nil, fmt.Errorf("db connection error: %w", err)
		}
	}

	{
		// connect to the mongodb
		ctx, done := context.WithTimeout(context.Background(), connectionTimeout)
		defer done()

		if err := client.Ping(ctx, readpref.Primary()); err != nil {
			return nil, fmt.Errorf("db primary ping error: %w", err)
		}
	}

	return client.Database(cfg.Name), nil
}
