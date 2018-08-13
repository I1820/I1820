/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 09-02-2018
 * |
 * | File Name:     query.go
 * +===============================================
 */

package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/mongodb/mongo-go-driver/bson"
)

// QueriesResource handles useful queries on database
type QueriesResource struct{}

type listResp struct {
	ID    string `json:"id" bson:"_id"`
	Total int    `json:"total"`
}

// List lists things and count of their data in database.
// This function is mapped to the path
// GET /queries/list
func (q QueriesResource) List(c buffalo.Context) error {
	var results = make([]listResp, 0)

	cur, err := db.Collection("data").Aggregate(c, bson.NewArray(
		bson.VC.DocumentFromElements(
			bson.EC.SubDocumentFromElements("$group",
				bson.EC.String("_id", "$thingid"),
				bson.EC.SubDocumentFromElements("total", bson.EC.Int32("$sum", 1)),
			),
		),
	))
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	for cur.Next(c) {
		var result listResp

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
