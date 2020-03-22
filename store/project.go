package store

import (
	"context"
	"fmt"

	"github.com/I1820/I1820/model"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func (ps Project) Get(ctx context.Context, name string) (model.Project, error) {
	var p model.Project

	dr := ps.DB.Collection(ProjectCollection).FindOne(ctx, bson.M{
		"name": name,
	})

	if err := dr.Decode(&p); err != nil {
		return p, err
	}

	return p, nil
}

func (ps Project) Delete(ctx context.Context, name string) error {
	// remove project entity from database
	if _, err := ps.DB.Collection(ProjectCollection).DeleteOne(ctx, bson.M{
		"name": name,
	}); err != nil {
		return err
	}

	if err := ps.DB.Collection(fmt.Sprintf("%s.logs.%s", ProjectCollection, name)).Drop(ctx); err != nil {
		logrus.Errorf("Log collection deletion failed %s", err)
	}

	return nil
}
