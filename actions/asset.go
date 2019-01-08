/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 20-09-2018
 * |
 * | File Name:     asset.go
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

// AssetsResource manages existing assets
type AssetsResource struct {
	buffalo.Resource
}

// asset request payload
type assetReq struct {
	Name  string `json:"name" validate:"alphanum,required"`
	Title string `json:"title" validate:"required"`
	Type  string `json:"type" validate:"required,oneof=boolean number string array object location"`
	Kind  string `json:"kind" validate:"required,oneof=sensor actuator"`
}

// List gets all assets of a given thing. This function is mapped to the path
// GET /projects/{project_id}/things/{thing_id}/assets
func (v AssetsResource) List(c buffalo.Context) error {
	thingID := c.Param("thing_id")
	projectID := c.Param("project_id")
	var t types.Thing

	dr := db.Collection("things").FindOne(c, bson.NewDocument(
		bson.EC.String("_id", thingID),
		bson.EC.String("project", projectID),
	))
	if err := dr.Decode(&t); err != nil {
		if err == mgo.ErrNoDocuments {
			return c.Error(http.StatusNotFound, fmt.Errorf("Thing %s not found", thingID))
		}
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, r.JSON(t.Assets))
}

// Create adds an asset to the DB and its thing. This function is mapped to the
// path POST /projects/{project_id}/things/{thing_id}/assets
func (v AssetsResource) Create(c buffalo.Context) error {
	thingID := c.Param("thing_id")
	projectID := c.Param("project_id")

	var rq assetReq
	if err := c.Bind(&rq); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}
	if err := validate.Struct(rq); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	a := types.Asset{
		Title: rq.Title,
		Type:  rq.Type,
		Kind:  rq.Kind,
	}

	dr := db.Collection("things").FindOneAndUpdate(c, bson.NewDocument(
		bson.EC.String("_id", thingID),
		bson.EC.String("project", projectID),
	), bson.NewDocument(
		bson.EC.SubDocumentFromElements("$set", bson.EC.Interface(fmt.Sprintf("assets.%s", rq.Name), a)),
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

// Show gets the data for a given asset. This function is mapped to
// the path GET /projects/{project_id}/things/{thing_id}/assets/{asset_name}
func (v AssetsResource) Show(c buffalo.Context) error {
	thingID := c.Param("thing_id")
	projectID := c.Param("project_id")
	assetName := c.Param("asset_id")

	var t types.Thing

	dr := db.Collection("things").FindOne(c, bson.NewDocument(
		bson.EC.String("_id", thingID),
		bson.EC.String("project", projectID),
	))

	if err := dr.Decode(&t); err != nil {
		if err == mgo.ErrNoDocuments {
			return c.Error(http.StatusNotFound, fmt.Errorf("Thing %s not found", thingID))
		}
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, r.JSON(t.Assets[assetName]))
}

// Destroy deletes an asset from the DB and its thing. This function is mapped
// to the path DELETE /projects/{project_id}/things/{thing_id}/assets/{asset_name}
func (v AssetsResource) Destroy(c buffalo.Context) error {
	thingID := c.Param("thing_id")
	projectID := c.Param("project_id")
	assetName := c.Param("asset_id")

	dr := db.Collection("things").FindOneAndUpdate(c, bson.NewDocument(
		bson.EC.String("_id", thingID),
		bson.EC.String("project", projectID),
	), bson.NewDocument(
		bson.EC.SubDocumentFromElements("$unset", bson.EC.String(fmt.Sprintf("assets.%s", assetName), "")),
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
