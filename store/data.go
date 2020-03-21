package store

import (
	"context"

	"github.com/I1820/I1820/model"
	"go.mongodb.org/mongo-driver/mongo"
)

// DataCollection is mongodb collection name for data
const DataCollection = "data"

// Data store handles database communication for data elements
type Data struct {
	DB *mongo.Database
}

// New creates new data store
func New(db *mongo.Database) *Data {
	return &Data{
		DB: db,
	}
}

// Insert given instance of data into database
func (d *Data) Insert(ctx context.Context, i model.Data) error {
	if _, err := d.DB.Collection(DataCollection).InsertOne(ctx, i); err != nil {
		return err
	}

	return nil
}
