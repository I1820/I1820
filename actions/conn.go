/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 24-09-2018
 * |
 * | File Name:     conn.go
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

// ConnectivitiesResource manages existing connectivities
type ConnectivitiesResource struct {
	buffalo.Resource
}

// connectivity request payload
type connectivityReq struct {
	Name string      `json:"name" validate:"alphanum,required"`
	Info interface{} `json:"info" validate:"required"`
}

// List gets all connectivities of a given thing. This function is mapped to the path
// GET /things/{thing_id}/connectivities
func (v ConnectivitiesResource) List(c buffalo.Context) error {
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

	return c.Render(http.StatusOK, r.JSON(t.Connectivities))
}

// Create adds a connectivity to the DB and its thing. This function is mapped to the
// path POST /things/{thing_id}/connectivities
func (v ConnectivitiesResource) Create(c buffalo.Context) error {
	thingID := c.Param("thing_id")

	var rq connectivityReq
	if err := c.Bind(&rq); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}
	if err := validate.Struct(rq); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	dr := db.Collection("things").FindOneAndUpdate(c, bson.NewDocument(
		bson.EC.String("_id", thingID),
	), bson.NewDocument(
		bson.EC.SubDocumentFromElements("$set", bson.EC.Interface(fmt.Sprintf("connectivities.%s", rq.Name), rq.Info)),
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

// Show gets the data for given connectivity. This function is mapped to
// the path GET /things/{thing_id}/connectivities/{connectivity_name}
func (v ConnectivitiesResource) Show(c buffalo.Context) error {
	thingID := c.Param("thing_id")
	connectivityName := c.Param("connectivity_id")

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

	return c.Render(http.StatusOK, r.JSON(t.Connectivities[connectivityName]))
}

// Destroy deletes a connectivity from the DB and its thing. This function is mapped
// to the path DELETE /things/{thing_id}/connectivities/{connectivity_name}
func (v ConnectivitiesResource) Destroy(c buffalo.Context) error {
	thingID := c.Param("thing_id")
	connectivityName := c.Param("connectivity_id")

	dr := db.Collection("things").FindOneAndUpdate(c, bson.NewDocument(
		bson.EC.String("_id", thingID),
	), bson.NewDocument(
		bson.EC.SubDocumentFromElements("$unset", bson.EC.String(fmt.Sprintf("connectivities.%s", connectivityName), "")),
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
