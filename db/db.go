package db

import (
	"context"
	"fmt"
	"time"

	"github.com/I1820/I1820/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionTimeout = 10 * time.Second

// New creates a new mongodb connection and tests it
func New(cfg config.Database) (*mongo.Database, error) {
	// create mongodb connection
	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.URL))
	if err != nil {
		return nil, fmt.Errorf("db new client error: %w", err)
	}

	// connect to the mongodb
	ctx, done := context.WithTimeout(context.Background(), connectionTimeout)
	defer done()

	if err := client.Connect(ctx); err != nil {
		return nil, fmt.Errorf("db connection error: %w", err)
	}

	return client.Database(cfg.Name), nil
}
