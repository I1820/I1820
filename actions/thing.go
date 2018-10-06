/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 07-07-2018
 * |
 * | File Name:     thing.go
 * +===============================================
 */

package actions

import (
	"fmt"
	"net/http"

	"github.com/I1820/pm/models"
	"github.com/I1820/types"
	"github.com/gobuffalo/buffalo"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	mgo "github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"github.com/mongodb/mongo-go-driver/mongo/mongoopt"
	"github.com/segmentio/ksuid"
)

// ThingsResource manages existing things
type ThingsResource struct {
	buffalo.Resource
}

// thing request payload
type thingReq struct {
	Name  string `json:"name" validate:"required"`
	Model string `json:"model"`
}

// List gets all things. This function is mapped to the path
// GET /projects/{project_id}/things
func (v ThingsResource) List(c buffalo.Context) error {
	projectID := c.Param("project_id")

	results := make([]types.Thing, 0)

	cur, err := db.Collection("things").Find(c, bson.NewDocument(
		bson.EC.String("project", projectID),
	))
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	for cur.Next(c) {
		var result types.Thing

		if err := cur.Decode(&result); err != nil {
			return c.Error(http.StatusInternalServerError, err)
		}

		results = append(results, result)
	}
	if err := cur.Close(c); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, r.JSON(results))
}

// Create adds a thing to the DB and its project. This function is mapped to the
// path POST /projects/{project_id}/things
func (v ThingsResource) Create(c buffalo.Context) error {
	projectID := c.Param("project_id")

	var rq thingReq
	if err := c.Bind(&rq); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	if err := validate.Struct(rq); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	// read more about thing model in I1820 platform website
	model := "generic"
	if rq.Model != "" {
		model = rq.Model
	}

	// check project existence
	dr := db.Collection("projects").FindOne(c, bson.NewDocument(
		bson.EC.String("_id", projectID),
	))

	var p models.Project
	if err := dr.Decode(&p); err != nil {
		if err == mgo.ErrNoDocuments {
			return c.Error(http.StatusNotFound, fmt.Errorf("Project %s not found", projectID))
		}
		return c.Error(http.StatusInternalServerError, err)
	}

	t := types.Thing{
		ID:             objectid.New().Hex(),
		Name:           rq.Name,
		Model:          model,
		Status:         true,
		Tokens:         []string{ksuid.New().String()},
		Assets:         make(map[string]types.Asset),
		Connectivities: make(map[string]interface{}),
		Tags:           make([]string, 0),

		Project: projectID,
	}

	if _, err := db.Collection("things").InsertOne(c, t); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	// create data collection with following format
	// data.project_id.thing_id
	cd := db.Collection(fmt.Sprintf("data.%s.%s", projectID, t.ID))
	if _, err := cd.Indexes().CreateMany(
		c,
		[]mgo.IndexModel{
			mgo.IndexModel{
				Keys: bson.NewDocument(
					bson.EC.Int32("at", 1),
				),
			},
			mgo.IndexModel{
				Keys: bson.NewDocument(
					bson.EC.Int32("asset", 1),
				),
			},
		},
	); err != nil {
		// this error should not happen but in case of it happens you can ignore it safely.
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, r.JSON(t))
}

// Show gets the data for one thing. This function is mapped to
// the path GET /projects/{project_id}/things/{thing_id}
func (v ThingsResource) Show(c buffalo.Context) error {
	projectID := c.Param("project_id")
	id := c.Param("thing_id")

	var t types.Thing

	dr := db.Collection("things").FindOne(c, bson.NewDocument(
		bson.EC.Boolean("status", true),
		bson.EC.String("_id", id),
		bson.EC.String("project", projectID),
	))

	if err := dr.Decode(&t); err != nil {
		if err == mgo.ErrNoDocuments {
			return c.Error(http.StatusNotFound, fmt.Errorf("Thing %s not found", id))
		}
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, r.JSON(t))
}

// Destroy deletes a thing from the DB and its project. This function is mapped
// to the path DELETE /projects/{project_id}/things/{thing_id}
func (v ThingsResource) Destroy(c buffalo.Context) error {
	projectID := c.Param("project_id")
	id := c.Param("thing_id")

	if _, err := db.Collection("things").DeleteOne(c, bson.NewDocument(
		bson.EC.String("_id", id),
		bson.EC.String("project", projectID),
	)); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, r.JSON(true))
}

// Activation activates/deactivates thing. This function is mapped
// to the path GET /projects/{project_id}/things/{thing_id}/{t:(?:activate|deactivate)}
func (v ThingsResource) Activation(c buffalo.Context) error {
	id := c.Param("thing_id")
	projectID := c.Param("project_id")

	status := false
	if c.Param("t") == "activate" {
		status = true
	}

	dr := db.Collection("things").FindOneAndUpdate(c, bson.NewDocument(
		bson.EC.String("_id", id),
		bson.EC.String("project", projectID),
	), bson.NewDocument(
		bson.EC.SubDocumentFromElements("$set", bson.EC.Boolean("status", status)),
	), findopt.ReturnDocument(mongoopt.After))

	var t types.Thing

	if err := dr.Decode(&t); err != nil {
		if err == mgo.ErrNoDocuments {
			return c.Error(http.StatusNotFound, fmt.Errorf("Thing %s not found", id))
		}
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, r.JSON(t))
}

// CreateToken creates new token for given device
func (v ThingsResource) CreateToken(c buffalo.Context) error {
	return nil
}

// RemoveToken removes token from given device
func (v ThingsResource) RemoveToken(c buffalo.Context) error {
	return nil
}
