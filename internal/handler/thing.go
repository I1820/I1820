package handler

import (
	"fmt"
	"net/http"

	"github.com/I1820/I1820/internal/model"
	"github.com/I1820/I1820/internal/request"
	"github.com/I1820/I1820/internal/store"
	"github.com/labstack/echo/v4"
)

// Things handles existing things.
type Things struct {
	Store store.Thing
}

// List gets all things. This function is mapped to the path
// GET /projects/{project_id}/things.
func (v Things) List(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	projectID := c.Param("project_id")

	results, err := v.Store.GetByProjectID(ctx, projectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, results)
}

// Create adds a thing to the DB and its project. This function is mapped to the
// path POST /projects/{project_id}/things.
func (v Things) Create(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	projectID := c.Param("project_id")

	var rq request.Thing
	if err := c.Bind(&rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := rq.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// read more about thing model in I1820 platform website
	if rq.Model == "" {
		rq.Model = "generic"
	}

	// there is no check for project existence!
	// but sjd-backend checks it

	t := model.Thing{
		Name:   rq.Name,
		Model:  rq.Model,
		Status: true,

		Project: projectID,
	}

	// set thing location if it is provided by user
	// otherwise it would be 0, 0
	/*
		t.Location.Type = "Point"
		t.Location.Coordinates = []float64{rq.Location.Longitude, rq.Location.Latitude}
	*/
	if err := v.Store.Create(ctx, t); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, t)
}

// Show gets the data for one thing. This function is mapped to
// the path GET /things/{thing_id}.
func (v Things) Show(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	id := c.Param("thing_id")

	t, err := v.Store.GetByName(ctx, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if t.Name == "" {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("thing %s not found", id))
	}

	return c.JSON(http.StatusOK, t)
}

// Update updates a thing information includes name, model and location. Please note that you must
// provide them all in update request even if you do not want to change it.
// This function is mapped to the path PUT /things/{thing_id}.
func (v Things) Update(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	id := c.Param("thing_id")

	var rq request.Thing
	if err := c.Bind(&rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := rq.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// read more about thing model in I1820 platform website
	if rq.Model == "" {
		rq.Model = "generic"
	}

	t, err := v.Store.Update(ctx, id, &rq.Model, nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if t.Name == "" {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("thing %s not found", id))
	}

	return c.JSON(http.StatusOK, t)
}

// Destroy deletes a thing from the DB and its project. This function is mapped
// to the path DELETE /things/{thing_id}.
func (v Things) Destroy(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	id := c.Param("thing_id")

	if err := v.Store.Remove(ctx, id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, true)
}

// Activation activates/deactivates thing. This function is mapped
// to the path GET /things/{thing_id}/{t:(?:activate|deactivate)}.
func (v Things) Activation(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	id := c.Param("thing_id")

	status := false
	if c.Param("t") == "activate" {
		status = true
	}

	t, err := v.Store.Update(ctx, id, nil, &status)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if t.Name == "" {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("thing %s not found", id))
	}

	return c.JSON(http.StatusOK, t)
}

// Register registers the routes of things handler on given echo group.
func (v Things) Register(g *echo.Group) {
	pg := g.Group("/projects/:project_id")
	{
		pg.GET("/things", v.List)
		pg.POST("/things", v.Create)
	}
	g.DELETE("/things/:thing_id", v.Destroy)
	g.GET("/things/:thing_id", v.Show)
	g.PUT("/things/:thing_id", v.Update)
	g.GET("/things/:thing_id/:t:(?:activate|deactivate)", v.Activation)
}
