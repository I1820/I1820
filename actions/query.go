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

	"github.com/I1820/types"
	"github.com/labstack/echo/v4"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// QueriesHandler handles useful queries on database
type QueriesHandler struct {
	db *mongo.Database
}

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

type recentlyReq struct {
	Asset  string `json:"asset"`
	Limit  int64  `json:"limit"`
	Offset int64  `json:"offset"`
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
func (q QueriesHandler) List(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	thingID := c.Param("thing_id")
	projectID := c.Param("project_id")

	var results = make([]listResp, 0)

	cur, err := q.db.Collection(fmt.Sprintf("data.%s.%s", projectID, thingID)).Aggregate(ctx, primitive.A{
		primitive.M{
			"$group": primitive.M{
				"_id":   "$asset",
				"total": primitive.M{"$sum": 1},
			},
		},
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	for cur.Next(ctx) {
		var result listResp

		if err := cur.Decode(&result); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		results = append(results, result)
	}
	if err := cur.Close(ctx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, results)
}

// Recently fetches given asset recent's data from database
// by default it fetches last 5 record.
// This function is mapped to the path
// POST projects/{project_id}/things/{thing_id}/queries/recently
func (q QueriesHandler) Recently(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	var req recentlyReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	limit := req.Limit
	if limit == 0 {
		limit = 5
	}
	offset := req.Offset
	assetName := req.Asset

	thingID := c.Param("thing_id")
	projectID := c.Param("project_id")

	cur, err := q.db.Collection(fmt.Sprintf("data.%s.%s", projectID, thingID)).Aggregate(ctx, primitive.A{
		primitive.M{
			"$match": primitive.M{
				"asset": assetName,
			},
		},
		primitive.M{
			"$sort": primitive.M{
				"at": -1,
			},
		},
		primitive.M{
			"$skip": offset,
		},
		primitive.M{
			"$limit": limit,
		},
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	results := make([]types.State, 0)
	for cur.Next(ctx) {
		var result types.State

		if err := cur.Decode(&result); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		results = append(results, result)
	}
	if err := cur.Close(ctx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, results)
}

// PartialFetch fetches data with windowing
// This function is mapped to the path
// POST projects/{project_id}/things/{thing_id}/queries/pfetch
// please consider that this query only works on numbers
func (q QueriesHandler) PartialFetch(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	var req fetchReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
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

	cur, err := q.db.Collection(fmt.Sprintf("data.%s.%s", projectID, thingID)).Aggregate(ctx, primitive.A{
		primitive.M{ // match phase
			"$match": primitive.M{
				"asset":        assetName,
				"value.number": primitive.M{"$exists": true},
				"at": primitive.M{
					"$gt": req.Range.From,
					"$lt": req.Range.To,
				},
			},
		},
		primitive.M{ // group phase
			"$group": primitive.M{
				"_id": primitive.M{
					"asset": "$asset",
					"cluster": primitive.M{
						"$floor": primitive.M{
							"$divide": primitive.A{
								primitive.M{
									"$subtract": primitive.A{
										"$at",
										0,
									},
								},
								cs,
							},
						},
					},
				},
				"count": primitive.M{"$sum": 1},
				"data":  primitive.M{"$avg": "$value.number"},
			},
		},
		primitive.M{ // add fields phase
			"$addFields": primitive.M{
				"since": primitive.M{
					"$add": primitive.A{
						0,
						primitive.M{
							"$multiply": primitive.A{
								"$_id.cluster",
								cs,
							},
						},
					},
				},
				"until": primitive.M{
					"$add": primitive.A{
						0,
						cs,
						primitive.M{
							"$multiply": primitive.A{
								"$_id.cluster",
								cs,
							},
						},
					},
				},
			},
		},
		primitive.M{ // sort phase
			"$sort": primitive.M{
				"since": -1,
			},
		},
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	results := make([]pfetchResp, 0)
	for cur.Next(ctx) {
		var result pfetchResp

		if err := cur.Decode(&result); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		results = append(results, result)
	}
	if err := cur.Close(ctx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, results)
}

// Fetch fetches given assets data in given time range from database.
// please consider that this fuction returns data in ascending time order.
// This function is mapped to the path
// POST projects/{project_id}/things/{thing_id}/queries/fetch
func (q QueriesHandler) Fetch(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	var req fetchReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	thingID := c.Param("thing_id")
	projectID := c.Param("project_id")

	assetName := req.Target
	assetType := req.Type

	cur, err := q.db.Collection(fmt.Sprintf("data.%s.%s", projectID, thingID)).Aggregate(ctx, primitive.A{
		primitive.M{
			"$match": primitive.M{ // find states of given asset that have given type
				"asset":                            assetName,
				fmt.Sprintf("value.%s", assetType): primitive.M{"$exists": true},
				"at": primitive.M{
					"$gt": req.Range.From,
					"$lt": req.Range.To,
				},
			},
		},
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	results := make([]types.State, 0)
	for cur.Next(ctx) {
		var result types.State

		if err := cur.Decode(&result); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		results = append(results, result)
	}
	if err := cur.Close(ctx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, results)
}
