/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 09-02-2018
 * |
 * | File Name:     query.go
 * +===============================================
 */

package handler

import (
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/I1820/I1820/request"
	"github.com/I1820/I1820/store"
	"github.com/labstack/echo/v4"
)

// Queries handles useful queries on database
type Queries struct {
	Store store.Data
}

// Register registers the routes of queries handler on given echo group
func (q Queries) Register(g *echo.Group) {
	g.GET("/queries/projects/:project_id/list", q.List)
	g.GET("/queries/things/:thing_id/fetch", q.FetchSingle)
	g.POST("/queries/fetch", q.Fetch)
}

// List lists things and count of their data in database.
// This function is mapped to the path
// GET /projects/{project_id}/things/{thing_id}/queries/list
func (q Queries) List(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	projectID := c.Param("project_id")

	results, err := q.Store.PerProjectCount(ctx, projectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, results)
}

// FetchSingle fetches the given thing data in given time range from database.
// please consider that this function returns data in ascending time order.
// This function is mapped to the path
// GET /queries/things/thing_id/fetch
func (q Queries) FetchSingle(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	id := c.Param("thing_id")

	since, err := strconv.ParseInt(c.QueryParam("since"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	until, err := strconv.ParseInt(c.QueryParam("until"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	results, err := q.Store.Fetch(ctx, since, until, 0, math.MaxInt64, []string{id})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, results)
}

// Fetch fetches given things data in given time range from database.
// please consider that this function returns data in ascending time order.
// This function is mapped to the path
// POST /queries/fetch
func (q Queries) Fetch(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	var req request.Fetch
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if req.Until == 0 {
		req.Until = time.Now().Unix()
	}

	if req.Limit == 0 {
		req.Limit = math.MaxInt64
	}

	results, err := q.Store.Fetch(ctx, req.Since, req.Until, req.Offset, req.Limit, req.ThingIDs)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, results)
}
