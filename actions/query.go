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
	"math"
	"net/http"
	"time"

	"github.com/I1820/types"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	ThingIDs []string `json:"thing_ids" validate:"required"`
	Since    int64    `json:"since" validate:"required"`
	Until    int64    `json:"until"`
	Limit    int64    `json:"limit"`
	Offset   int64    `json:"offset"`
}

// List lists assets and count of their data in database.
// This function is mapped to the path
// GET /projects/{project_id}/things/{thing_id}/qeuries/list
func (q QueriesHandler) List(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	projectID := c.Param("project_id")

	var results = make([]listResp, 0)

	cur, err := q.db.Collection("data").Aggregate(ctx, bson.A{
		bson.M{
			"$match": bson.M{
				"project": projectID,
			},
		},
		bson.M{
			"$group": bson.M{
				"_id":   "$thingid",
				"total": bson.M{"$sum": 1},
			},
		},
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	for cur.Next(ctx) {
		var result listResp

		if err := cur.Decode(&result); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		results = append(results, result)
	}
	if err := cur.Close(ctx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, results)
}

// Fetch fetches given things data in given time range from database.
// please consider that this fuction returns data in ascending time order.
// This function is mapped to the path
// POST /queries/fetch
func (q QueriesHandler) Fetch(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	var req fetchReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if req.Until == 0 {
		req.Until = time.Now().Unix()
	}

	if req.Limit == 0 {
		req.Limit = math.MaxInt64
	}

	cur, err := q.db.Collection("data").Aggregate(ctx, bson.A{
		bson.M{
			"$match": bson.M{
				"thingid": bson.M{
					"$in": req.ThingIDs,
				},
				"timestamp": bson.M{
					"$gt": time.Unix(req.Since, 0),
					"$lt": time.Unix(req.Until, 0),
				},
			},
		},
		bson.M{
			"$sort": bson.M{
				"timestamp": -1,
			},
		},
		bson.M{
			"$skip": req.Offset,
		},
		bson.M{
			"$limit": req.Limit,
		},
	}, options.Aggregate().SetAllowDiskUse(true))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	results := make([]types.Data, 0)
	for cur.Next(ctx) {
		var result types.Data

		if err := cur.Decode(&result); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		results = append(results, result)

		fmt.Printf("%+v\n", result)
	}
	if err := cur.Close(ctx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, results)
}
