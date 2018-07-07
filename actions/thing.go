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
	"net/http"

	"github.com/aiotrc/pm/models"
	"github.com/gobuffalo/buffalo"
	"github.com/mongodb/mongo-go-driver/bson"
)

// ThingsResource manages existing things
type ThingsResource struct {
	buffalo.Resource
}

// thing request payload
type thingReq struct {
	Name string `json:"name" binding:"required"`
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
