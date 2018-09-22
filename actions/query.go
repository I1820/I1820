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
	"fmt"
	"net/http"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/mongodb/mongo-go-driver/bson"
)

// QueriesResource handles useful queries on database
type QueriesResource struct{}

type listResp struct {
	ID    string `json:"id" bson:"_id"`
	Total int    `json:"total"`
}

type fetchReq struct {
	Range struct {
		From time.Time `json:"from"`
		To   time.Time `json:"to"`
	} `json:"range"`
	IntervalMs int64 `json:"intervalMs"`
	Targets    []struct {
		Target string `json:"target"`
		RefID  string `json:"refId"`
	} `json:"targets"`
}

type fetchResp struct {
	Target     string
	Datapoints []struct {
		Metric    float64
		Timestamp int64
	}
}

// List lists assets and count of their data in database.
// This function is mapped to the path
// GET /projects/{project_id}/things/{thing_id}/qeuries/list
func (q QueriesResource) List(c buffalo.Context) error {
	thingID := c.Param("thing_id")
	projectID := c.Param("project_id")

	var results = make([]listResp, 0)

	cur, err := db.Collection(fmt.Sprintf("data.%s.%s", projectID, thingID)).Aggregate(c, bson.NewArray(
		bson.VC.DocumentFromElements(
			bson.EC.SubDocumentFromElements("$group",
				bson.EC.String("_id", "$asset"),
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

// Fetch fetches given keys data from database.
// it works in given range based on given intervals
// This function is mapped to the path
// POST /queries/{}/fetch
func (q QueriesResource) Fetch(c buffalo.Context) error {
	var req fetchReq
	if err := c.Bind(&req); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}
	fmt.Println(req)

	cur, err := db.Collection("data").Aggregate(c, bson.NewArray(
		bson.VC.DocumentFromElements(
			bson.EC.SubDocumentFromElements("$match",
				bson.EC.String("project", c.Param("project_id")),
				bson.EC.SubDocumentFromElements("data", bson.EC.Null("$ne")),
				bson.EC.SubDocumentFromElements("timestamp",
					bson.EC.Time("$gt", req.Range.From),
					bson.EC.Time("$lt", req.Range.To),
				),
			),
		),
		bson.VC.DocumentFromElements(
			bson.EC.SubDocumentFromElements("$project",
				bson.EC.SubDocumentFromElements("targets", bson.EC.String("$objectToArray", "$data")),
			),
		),
		bson.VC.DocumentFromElements(bson.EC.String("$unwind", "$targets")),
	))
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	for cur.Next(c) {
		var result struct {
			Targets struct {
				K interface{}
				V interface{}
			}
		}

		if err := cur.Decode(&result); err != nil {
			return c.Error(http.StatusInternalServerError, err)
		}

		fmt.Println(result)
	}
	if err := cur.Close(c); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	var results = make([]fetchResp, 0)

	return c.Render(http.StatusOK, r.JSON(results))
}
