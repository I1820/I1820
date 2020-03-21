package store

import (
	"context"
	"time"

	"github.com/I1820/I1820/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// PerProjectCount counts number of records for given project per things
// it returns a map between thing identification and number of records for that thing
func (d Data) PerProjectCount(ctx context.Context, projectID string) (map[string]int, error) {
	cur, err := d.DB.Collection(DataCollection).Aggregate(ctx, bson.A{
		bson.M{
			"$match": bson.M{
				"project": projectID,
			},
		},
		bson.M{
			"$group": bson.M{
				"_id":   "$thingid",
				"total": bson.M{"$sum": 1},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	results := make(map[string]int)

	for cur.Next(ctx) {
		var result struct {
			ID    string `bson:"_id"`
			Total int    `bson:"total"`
		}

		if err := cur.Decode(&result); err != nil {
			return nil, err
		}

		results[result.ID] = result.Total
	}

	if err := cur.Close(ctx); err != nil {
		return nil, err
	}

	return results, err
}

// Fetch fetches given things data in given time range and sorts it.
func (d Data) Fetch(ctx context.Context, since, until, offset, limit int64, ids []string) ([]model.Data, error) {
	cur, err := d.DB.Collection(DataCollection).Aggregate(ctx, bson.A{
		bson.M{
			"$match": bson.M{
				"thingid": bson.M{
					"$in": ids,
				},
				"timestamp": bson.M{
					"$gt": time.Unix(since, 0),
					"$lt": time.Unix(until, 0),
				},
			},
		},
		bson.M{
			"$sort": bson.M{
				"timestamp": -1,
			},
		},
		bson.M{
			"$skip": offset,
		},
		bson.M{
			"$limit": limit,
		},
	}, options.Aggregate().SetAllowDiskUse(true))
	if err != nil {
		return nil, err
	}

	results := make([]model.Data, 0)

	for cur.Next(ctx) {
		var result model.Data

		if err := cur.Decode(&result); err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	if err := cur.Close(ctx); err != nil {
		return nil, err
	}

	return results, nil
}
