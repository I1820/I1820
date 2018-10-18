/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 18-10-2018
 * |
 * | File Name:     tag.go
 * +===============================================
 */

package actions

import (
	"fmt"
	"net/http"

	"github.com/I1820/types"
	"github.com/gobuffalo/buffalo"
	"github.com/mongodb/mongo-go-driver/bson"
	mgo "github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"github.com/mongodb/mongo-go-driver/mongo/mongoopt"
)

type TagsResource struct{}

// tagReq contains of user given tag array.
type tagReq struct {
	Tags []string `json:"tags" validate:"required"`
}

// Create replaces tag array of the given thing with the newly given tag array.
// This function is mapped to the path POST /things/{thing_id}/tags
func (v TagsResource) Create(c buffalo.Context) error {
	thingID := c.Param("thing_id")

	var rq tagReq
	if err := c.Bind(&rq); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}
	if err := validate.Struct(rq); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	dr := db.Collection("things").FindOneAndUpdate(c, bson.NewDocument(
		bson.EC.String("_id", thingID),
	), bson.NewDocument(
		bson.EC.SubDocumentFromElements("$set", bson.EC.Interface("tags", rq.Tags)),
	), findopt.ReturnDocument(mongoopt.After))

	var t types.Thing

	if err := dr.Decode(&t); err != nil {
		if err == mgo.ErrNoDocuments {
			return c.Error(http.StatusNotFound, fmt.Errorf("Thing %s not found", thingID))
		}
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, r.JSON(t))
}

// List returns the tag array of the given thing.
// This function is mapped to the path GET /things/{thing_id}/tags
func (v TagsResource) List(c buffalo.Context) error {
	thingID := c.Param("thing_id")
	var t types.Thing

	dr := db.Collection("things").FindOne(c, bson.NewDocument(
		bson.EC.String("_id", thingID),
	))
	if err := dr.Decode(&t); err != nil {
		if err == mgo.ErrNoDocuments {
			return c.Error(http.StatusNotFound, fmt.Errorf("Thing %s not found", thingID))
		}
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, r.JSON(t.Tags))
}
