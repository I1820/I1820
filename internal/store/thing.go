package store

import (
	"context"

	"github.com/I1820/I1820/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ThingCollection is mongodb collection name for things
const ThingCollection = "thing"

// Thing stores and retrieves things collection
type Thing struct {
	DB *mongo.Database
}

// GetByProjectID returns all things that are associated with given project identification
func (ts Thing) GetByProjectID(ctx context.Context, pid string) ([]model.Thing, error) {
	var things []model.Thing

	cur, err := ts.DB.Collection(ThingCollection).Find(ctx, bson.M{
		"project": pid,
	})
	if err != nil {
		return things, err
	}

	for cur.Next(ctx) {
		var thing model.Thing

		if err := cur.Decode(&thing); err != nil {
			return things, err
		}

		things = append(things, thing)
	}

	if err := cur.Close(ctx); err != nil {
		return things, err
	}

	return things, nil
}

// GetByName returns the thing that has given identification
func (ts Thing) GetByName(ctx context.Context, id string) (model.Thing, error) {
	var t model.Thing

	dr := ts.DB.Collection(ThingCollection).FindOne(ctx, bson.M{
		"status": true,
		"name":   id,
	})

	if err := dr.Decode(&t); err != nil {
		if err == mongo.ErrNoDocuments {
			return model.Thing{}, nil
		}

		return t, err
	}

	return t, nil
}

// Create creates the given thing in things collection
func (ts Thing) Create(ctx context.Context, t model.Thing) error {
	if _, err := ts.DB.Collection(ThingCollection).InsertOne(ctx, t); err != nil {
		return err
	}

	return nil
}

// Update update the given thing's model
func (ts Thing) Update(ctx context.Context, id string, m *string, s *bool) (model.Thing, error) {
	set := bson.M{}

	if m != nil {
		set["model"] = *m
	}

	if s != nil {
		set["status"] = *s
	}

	dr := ts.DB.Collection(ThingCollection).FindOneAndUpdate(ctx, bson.M{
		"name": id,
	}, bson.M{
		"$set": set,
	}, options.FindOneAndUpdate().SetReturnDocument(options.After))

	var t model.Thing

	if err := dr.Decode(&t); err != nil {
		if err == mongo.ErrNoDocuments {
			return model.Thing{}, nil
		}

		return t, err
	}

	return t, nil
}

// Remove removes the given thing from the things collection
func (ts Thing) Remove(ctx context.Context, id string) error {
	if _, err := ts.DB.Collection(ThingCollection).DeleteOne(ctx, bson.M{
		"name": id,
	}); err != nil {
		return err
	}

	return nil
}
