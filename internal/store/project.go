package store

import (
	"context"
	"fmt"

	"github.com/I1820/I1820/internal/model"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ProjectCollection is mongodb collection name for project
const ProjectCollection = "project"

// Project stores and retrieves projects collection
type Project struct {
	DB *mongo.Database
}

// List returns all registered projects
func (ps Project) List(ctx context.Context) ([]model.Project, error) {
	projects := make([]model.Project, 0)

	cur, err := ps.DB.Collection(ProjectCollection).Find(ctx, bson.M{})
	if err != nil {
		return projects, err
	}

	for cur.Next(ctx) {
		var p model.Project

		if err := cur.Decode(&p); err != nil {
			return projects, err
		}

		projects = append(projects, p)
	}

	if err := cur.Close(ctx); err != nil {
		return projects, err
	}

	return projects, err
}

func (ps Project) Create(ctx context.Context, p model.Project) error {
	if _, err := ps.DB.Collection(ProjectCollection).InsertOne(ctx, p); err != nil {
		return err
	}

	return nil
}

func (ps Project) Get(ctx context.Context, id string) (model.Project, error) {
	var p model.Project

	dr := ps.DB.Collection(ProjectCollection).FindOne(ctx, bson.M{
		"id": id,
	})

	if err := dr.Decode(&p); err != nil {
		return p, err
	}

	return p, nil
}

func (ps Project) Delete(ctx context.Context, id string) error {
	// remove project entity from database
	if _, err := ps.DB.Collection(ProjectCollection).DeleteOne(ctx, bson.M{
		"id": id,
	}); err != nil {
		return err
	}

	if err := ps.DB.Collection(fmt.Sprintf("%s.logs.%s", ProjectCollection, id)).Drop(ctx); err != nil {
		logrus.Errorf("Log collection deletion failed %s", err)
	}

	return nil
}

func (ps Project) Update(ctx context.Context, id string, fields map[string]interface{}) (model.Project, error) {
	var p model.Project

	dr := ps.DB.Collection(ProjectCollection).FindOneAndUpdate(ctx, bson.M{
		"id": id,
	}, bson.M{
		"$set": fields,
	}, options.FindOneAndUpdate().SetReturnDocument(options.After))

	if err := dr.Decode(&p); err != nil {
		return p, err
	}

	return p, nil
}

func (ps Project) Logs(ctx context.Context, name string, limit int) ([]model.ProjectLog, error) {
	var pls = make([]model.ProjectLog, 0)

	cur, err := ps.DB.Collection(fmt.Sprintf("%s.logs.%s", ProjectCollection, name)).Aggregate(ctx, bson.A{
		bson.M{
			"$sort": bson.M{
				"Time": -1,
			},
		},
		bson.M{
			"$limit": limit,
		},
	}, options.Aggregate().SetAllowDiskUse(true))
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var pl model.ProjectLog

		if err := cur.Decode(&pl); err != nil {
			return nil, err
		}

		pls = append(pls, pl)
	}

	if err := cur.Close(ctx); err != nil {
		return nil, err
	}

	return pls, nil
}
