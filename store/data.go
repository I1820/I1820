/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 18-07-2019
 * |
 * | File Name:     data.go
 * +===============================================
 */

package store

import (
	"context"
	"time"

	"github.com/I1820/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collection = "data"

// Data stores and retrieves data collection
type Data struct {
	DB *mongo.Database
}

// PerProjectCount counts number of records for given project per things
// it returns a map between thing identification and number of records for that thing
func (ds Data) PerProjectCount(ctx context.Context, projectID string) (map[string]int, error) {
	cur, err := ds.DB.Collection(collection).Aggregate(ctx, bson.A{
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
func (ds Data) Fetch(ctx context.Context, since, until, offset, limit int64, ids []string) ([]types.Data, error) {
	cur, err := ds.DB.Collection(collection).Aggregate(ctx, bson.A{
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

	results := make([]types.Data, 0)
	for cur.Next(ctx) {
		var result types.Data

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
