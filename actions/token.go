/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 17-11-2018
 * |
 * | File Name:     token.go
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
	"github.com/segmentio/ksuid"
)

// TokensResource manages existing assets
type TokensResource struct {
}

// Create creates new token for given device
// path GET /projects/{project_id}/things/{thing_id}/tokens
func (v TokensResource) Create(c buffalo.Context) error {
	projectID := c.Param("project_id")
	id := c.Param("thing_id")

	var t types.Thing

	dr := db.Collection("things").FindOneAndUpdate(c, bson.NewDocument(
		bson.EC.String("_id", id),
		bson.EC.String("project", projectID),
	), bson.NewDocument(
		bson.EC.SubDocumentFromElements("$addToSet", bson.EC.String("tokens", ksuid.New().String())),
	), findopt.ReturnDocument(mongoopt.After))

	if err := dr.Decode(&t); err != nil {
		if err == mgo.ErrNoDocuments {
			return c.Error(http.StatusNotFound, fmt.Errorf("Thing %s not found", id))
		}
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, r.JSON(t))
}

// Destroy removes token from given device
// path DELETE /projects/{project_id}/things/{thing_id}/tokens/{token}
func (v TokensResource) Destroy(c buffalo.Context) error {
	projectID := c.Param("project_id")
	id := c.Param("thing_id")
	token := c.Param("token")

	var t types.Thing

	dr := db.Collection("things").FindOneAndUpdate(c, bson.NewDocument(
		bson.EC.String("_id", id),
		bson.EC.String("project", projectID),
	), bson.NewDocument(
		bson.EC.SubDocumentFromElements("$pull", bson.EC.String("tokens", token)),
	), findopt.ReturnDocument(mongoopt.After))

	if err := dr.Decode(&t); err != nil {
		if err == mgo.ErrNoDocuments {
			return c.Error(http.StatusNotFound, fmt.Errorf("Thing %s not found", id))
		}
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, r.JSON(t))
}
