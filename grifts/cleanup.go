package grifts

import (
	"log"

	"github.com/aiotrc/pm/models"
	"github.com/gobuffalo/envy"
	grift "github.com/markbates/grift/grift"
	"github.com/mongodb/mongo-go-driver/bson"
	mgo "github.com/mongodb/mongo-go-driver/mongo"
)

var _ = grift.Desc("cleanup", "Removes existing projects")
var _ = grift.Add("cleanup", func(c *grift.Context) error {
	// Create mongodb connection
	url := envy.Get("DB_URL", "mongodb://127.0.0.1")
	client, err := mgo.NewClient(url)
	if err != nil {
		return err
	}
	if err := client.Connect(c); err != nil {
		return err
	}
	log.Printf("DB url: %s\n", url)

	db := client.Database("isrc")

	// Get all project
	ps := make([]models.Project, 0)

	cur, err := db.Collection("pm").Find(c, bson.NewDocument())
	if err != nil {
		return err
	}

	for cur.Next(c) {
		var p models.Project

		if err := cur.Decode(&p); err != nil {
			return err
		}

		ps = append(ps, p)
	}
	if err := cur.Close(c); err != nil {
		return err
	}

	log.Printf("Projects: %v", ps)

	// Remove all projects
	for _, p := range ps {
		if err := p.Runner.Remove(c); err != nil {
			return err
		}

		if _, err := db.Collection("pm").DeleteOne(c, bson.NewDocument(
			bson.EC.String("name", p.Name),
		)); err != nil {
			return err
		}

		log.Printf("Project %s was removed", p.Name)
	}

	return nil
})
