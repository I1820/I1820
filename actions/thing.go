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

	"github.com/aiotrc/pm/models"
	"github.com/gobuffalo/buffalo"
	"github.com/mongodb/mongo-go-driver/bson"
	mgo "github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"github.com/mongodb/mongo-go-driver/mongo/mongoopt"
)

// ThingsResource manages existing things
type ThingsResource struct {
	buffalo.Resource
}

// thing request payload
type thingReq struct {
	Name    string `json:"name" binding:"required"`
	Project string `json:"project" binding:"required"`
}

// List gets all things. This function is mapped to the path
// GET /things
func (v ThingsResource) List(c buffalo.Context) error {
	results := make([]models.Thing, 0)

	cur, err := db.Collection("pm").Aggregate(c, bson.NewArray(
		bson.VC.DocumentFromElements(
			bson.EC.String("$unwind", "$things"),
		),
		bson.VC.DocumentFromElements(
			bson.EC.SubDocumentFromElements("$replaceRoot", bson.EC.String("newRoot", "$things")),
		),
	))
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	for cur.Next(c) {
		var result models.Thing

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
// path POST /things
func (v ThingsResource) Create(c buffalo.Context) error {
	var rq thingReq
	if err := c.Bind(&rq); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	project := rq.Project

	t := models.Thing{
		ID:     rq.Name,
		Status: true,
	}

	dr := db.Collection("pm").FindOneAndUpdate(c, bson.NewDocument(
		bson.EC.String("name", project),
	), bson.NewDocument(
		bson.EC.SubDocumentFromElements("$addToSet", bson.EC.Interface("things", t)),
	), findopt.ReturnDocument(mongoopt.After))

	var p models.Project

	if err := dr.Decode(&p); err != nil {
		if err == mgo.ErrNoDocuments {
			return c.Error(http.StatusNotFound, fmt.Errorf("Project %s not found", project))
		}
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, r.JSON(p))
}

// Show gets the project for one thing. This function is mapped to
// the path GET /things/{thing_id}
func (v ThingsResource) Show(c buffalo.Context) error {
	name := c.Param("thing_id")

	var p models.Project

	dr := db.Collection("pm").FindOne(c, bson.NewDocument(
		bson.EC.Boolean("status", true),
		bson.EC.SubDocumentFromElements("things", bson.EC.SubDocumentFromElements("$elemMatch",
			bson.EC.String("id", name),
		)),
	))

	if err := dr.Decode(&p); err != nil {
		if err == mgo.ErrNoDocuments {
			return c.Error(http.StatusNotFound, fmt.Errorf("Thing %s not found", name))
		}
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, r.JSON(p))
}

// Destroy deletes a thing from the DB and its project. This function is mapped
// to the path DELETE /things/{thing_id}
func (v ThingsResource) Destroy(c buffalo.Context) error {
	name := c.Param("thing_id")

	dr := db.Collection("pm").FindOneAndUpdate(c, bson.NewDocument(), bson.NewDocument(
		bson.EC.SubDocumentFromElements("$pull", bson.EC.SubDocumentFromElements(
			"things", bson.EC.String("id", name)),
		),
	), findopt.ReturnDocument(mongoopt.After))

	var p models.Project

	if err := dr.Decode(&p); err != nil {
		if err == mgo.ErrNoDocuments {
			return c.Error(http.StatusNotFound, fmt.Errorf("Thing %s not found", name))
		}
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, r.JSON(p))
}

// Activation activates/deactivates thing. This function is mapped
// to the path GET /things/{thing_id}/{t:(?:activate|deactivate)}
func (v ThingsResource) Activation(c buffalo.Context) error {
	name := c.Param("thing_id")

	t := c.Param("t")
	status := false
	if t == "activate" {
		status = true
	}

	dr := db.Collection("pm").FindOneAndUpdate(c, bson.NewDocument(
		bson.EC.SubDocumentFromElements("things", bson.EC.SubDocumentFromElements(
			"$elemMatch", bson.EC.String("id", name),
		)),
	), bson.NewDocument(
		bson.EC.SubDocumentFromElements("$set", bson.EC.Boolean("things.$.status", status)),
	), findopt.ReturnDocument(mongoopt.After))

	var p models.Project

	if err := dr.Decode(&p); err != nil {
		if err == mgo.ErrNoDocuments {
			c.Error(http.StatusNotFound, fmt.Errorf("Thing %s not found", name))
		}
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, r.JSON(p))
}
