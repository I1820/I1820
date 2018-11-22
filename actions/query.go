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
	Type   string `json:"type"`
	Target string `json:"target"`
	Window struct {
		Size int64 `json:"size"`
	} `json:"window"`
}

type pfetchResp struct {
	ID struct {
		Asset   string `json:"asset" bson:"asset"`
		Cluster int64  `json:"cluster" bson:"cluster"`
	} `json:"id" bson:"_id"`
	Count int       `json:"count" bson:"count"`
	Data  float64   `json:"data" bson:"data"`
	Since time.Time `json:"since" bson:"since"`
	Until time.Time `json:"until" bson:"until"`
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

// PartialFetch fetches data with windowing
// This function is mapped to the path
// POST projects/{project_id}/things/{thing_id}/queries/pfetch
// please consider that this query only works on numbers
func (q QueriesResource) PartialFetch(c buffalo.Context) error {
	var req fetchReq
	if err := c.Bind(&req); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	thingID := c.Param("thing_id")
	projectID := c.Param("project_id")
	assetName := req.Target

	// set default window size
	if req.Window.Size == 0 {
		req.Window.Size = 200
	}

	// to - from / window size indicates each partition duration in milliseconds
	cs := int64(req.Range.To.Sub(req.Range.From).Seconds()*1000) / req.Window.Size
	if cs == 0 {
		cs++
	}

	cur, err := db.Collection(fmt.Sprintf("data.%s.%s", projectID, thingID)).Aggregate(c, bson.NewArray(
		bson.VC.DocumentFromElements( // match phase
			bson.EC.SubDocumentFromElements("$match",
				bson.EC.String("asset", assetName),
				bson.EC.SubDocumentFromElements("value.number", bson.EC.Boolean("$exists", true)),
				bson.EC.SubDocumentFromElements("at",
					bson.EC.Time("$gt", req.Range.From),
					bson.EC.Time("$lt", req.Range.To),
				),
			),
		),
		bson.VC.DocumentFromElements( // group phase
			bson.EC.SubDocumentFromElements("$group",
				bson.EC.SubDocumentFromElements("_id",
					bson.EC.String("asset", "$asset"),
					bson.EC.SubDocumentFromElements("cluster",
						bson.EC.SubDocumentFromElements("$floor",
							bson.EC.ArrayFromElements("$divide",
								bson.VC.DocumentFromElements(
									bson.EC.ArrayFromElements("$subtract",
										bson.VC.String("$at"),
										bson.VC.DateTime(0),
									),
								),
								bson.VC.Int64(cs),
							),
						),
					),
				),
				bson.EC.SubDocumentFromElements("count", bson.EC.Int32("$sum", 1)),
				bson.EC.SubDocumentFromElements("data", bson.EC.String("$avg", "$value.number")),
			),
		),
		bson.VC.DocumentFromElements( // add fields phase
			bson.EC.SubDocumentFromElements("$addFields",
				bson.EC.SubDocumentFromElements("since",
					bson.EC.ArrayFromElements("$add",
						bson.VC.DateTime(0),
						bson.VC.DocumentFromElements(
							bson.EC.ArrayFromElements("$multiply",
								bson.VC.String("$_id.cluster"),
								bson.VC.Int64(cs),
							),
						),
					),
				),
				bson.EC.SubDocumentFromElements("until",
					bson.EC.ArrayFromElements("$add",
						bson.VC.DateTime(0),
						bson.VC.Int64(cs),
						bson.VC.DocumentFromElements(
							bson.EC.ArrayFromElements("$multiply",
								bson.VC.String("$_id.cluster"),
								bson.VC.Int64(cs),
							),
						),
					),
				),
			),
		),
		bson.VC.DocumentFromElements( // sort phase
			bson.EC.SubDocumentFromElements("$sort",
				bson.EC.Int32("since", -1),
			),
		),
	))
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	results := make([]pfetchResp, 0)
	for cur.Next(c) {
		var result pfetchResp

		if err := cur.Decode(&result); err != nil {
			return c.Error(http.StatusInternalServerError, err)
		}

		results = append(results, result)
	}
	if err := cur.Close(c); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}
	fmt.Println(results)

	return c.Render(http.StatusOK, r.JSON(results))
}

// Fetch fetches given assets data in given time range from database.
// please consider that this fuction returns data in ascending time order.
// This function is mapped to the path
// POST projects/{project_id}/things/{thing_id}/queries/fetch
func (q QueriesResource) Fetch(c buffalo.Context) error {
	var req fetchReq
	if err := c.Bind(&req); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	thingID := c.Param("thing_id")
	projectID := c.Param("project_id")

	assetName := req.Target
	assetType := req.Type

	cur, err := db.Collection(fmt.Sprintf("data.%s.%s", projectID, thingID)).Aggregate(c, bson.NewArray(
		bson.VC.DocumentFromElements(
			bson.EC.SubDocumentFromElements("$match", // find states of given asset that have given type
				bson.EC.String("asset", assetName),
				bson.EC.SubDocumentFromElements(fmt.Sprintf("value.%s", assetType), bson.EC.Boolean("$exists", true)),
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
