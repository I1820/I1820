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
	"strconv"
	"time"

	"github.com/I1820/types"
	"github.com/gobuffalo/buffalo"
	"github.com/mongodb/mongo-go-driver/bson"
	mgo "github.com/mongodb/mongo-go-driver/mongo"
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
	Target string `json:"target"`
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

// Recently fetches given assets recently data from database
// by default it fetches last 5 record
// This function is mapped to the path
// POST projects/{project_id}/things/{thing_id}/assets/{asset_name}/queries/recently
func (q QueriesResource) Recently(c buffalo.Context) error {
	limit, err := strconv.ParseInt(c.Param("limit"), 10, 64)
	if err != nil {
		limit = 5
	}
	assetName := c.Param("asset_name")
	thingID := c.Param("thing_id")
	projectID := c.Param("project_id")

	cur, err := db.Collection(fmt.Sprintf("data.%s.%s", projectID, thingID)).Aggregate(c, bson.NewArray(
		bson.VC.DocumentFromElements(
			bson.EC.SubDocumentFromElements("$match",
				bson.EC.String("asset", assetName),
			),
		),
		bson.VC.DocumentFromElements(
			bson.EC.SubDocumentFromElements("$sort",
				bson.EC.Int32("at", -1),
			),
		),
		bson.VC.DocumentFromElements(
			bson.EC.Int64("$limit", limit),
		),
	))
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	results := make([]types.State, 0)
	for cur.Next(c) {
		var result types.State

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

// Fetch fetches given assets data in given time range from database.
// This function is mapped to the path
// POST projects/{project_id}/things/{thing_id}/queries/fetch
func (q QueriesResource) Fetch(c buffalo.Context) error {
	var req fetchReq
	if err := c.Bind(&req); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	thingID := c.Param("thing_id")
	projectID := c.Param("project_id")

	// find things by its id
	// then find given asset and its type
	var t types.Thing
	dr := db.Collection("things").FindOne(c, bson.NewDocument(
		bson.EC.String("_id", thingID),
	))
	if err := dr.Decode(&t); err != nil {
		if err == mgo.ErrNoDocuments {
			return c.Error(http.StatusInternalServerError, fmt.Errorf("Thing %s not found", thingID))
		}
		return c.Error(http.StatusInternalServerError, err)
	}
	assetName := req.Target
	assetType := t.Assets[assetName].Type
	if assetType == "" {
		assetType = "String"
	}

	fmt.Println(assetName, assetType)

	cur, err := db.Collection(fmt.Sprintf("data.%s.%s", projectID, thingID)).Aggregate(c, bson.NewArray(
		bson.VC.DocumentFromElements(
			bson.EC.SubDocumentFromElements("$match",
				bson.EC.String("asset", assetName),
				// bson.EC.SubDocumentFromElements(fmt.Sprintf("value.%s", assetType), bson.EC.Boolean("$exists", true)),
				bson.EC.SubDocumentFromElements("at",
					bson.EC.Time("$gt", req.Range.From),
					bson.EC.Time("$lt", req.Range.To),
				),
			),
		),
	))
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	results := make([]types.State, 0)
	for cur.Next(c) {
		var result types.State

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
