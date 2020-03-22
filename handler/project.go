package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/I1820/I1820/config"
	"github.com/I1820/I1820/model"
	"github.com/I1820/I1820/request"
	"github.com/I1820/I1820/runner"
	"github.com/I1820/I1820/store"
	"github.com/labstack/echo/v4"
)

// Projects manages existing projects
type Projects struct {
	Store   store.Project
	Manager *runner.Manager

	Config config.Runner
}

func (v Projects) Register(g *echo.Group) {
	g.POST("/projects", v.Create)
	g.DELETE("/projects/:project_id", v.Destroy)
	g.GET("/projects/:project_id", v.Show)
	g.PUT("/projects/:project_id", v.Update)
	g.GET("/projects/:project_id/logs", v.Logs)
	g.GET("/projects/:project_id/recreate", v.Recreate)
}

// List gets all projects. This function is mapped to the path
// GET /projects
func (v Projects) List(c echo.Context) error {
	ctx := c.Request().Context()

	projects, err := v.Store.List(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, projects)
}

// Create adds a project to the DB and creates its docker. This function is mapped to the
// path POST /projects
func (v Projects) Create(c echo.Context) error {
	ctx := c.Request().Context()

	var rq request.Project

	if err := c.Bind(&rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := rq.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	addr := fmt.Sprintf(
		"amqp://%s:%s@%s:%d/",
		v.Config.Rabbitmq.User,
		v.Config.Rabbitmq.Pass,
		v.Config.Rabbitmq.Host,
		v.Config.Rabbitmq.Port,
	)

	// predefined environment variables
	envs := []runner.Env{
		{Name: "DB_URL", Value: v.Config.Database.URL},
		{Name: "BROKER_URL", Value: addr},
		{Name: "OWNER", Value: rq.Owner},
	}

	// user-defined environment variables
	for envKey, envVal := range rq.Envs {
		envs = append(envs, runner.Env{Name: envKey, Value: envVal})
	}

	var p model.Project

	id := model.NewProjectID()

	// creates project entity with its docker (have fun :D)
	r, err := v.Manager.New(ctx, id, envs)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	p.Runner = r

	// sets other properties of the project
	p.ID = id
	p.Description = rq.Description

	// converts request location to GeoJSON format
	p.Perimeter.Type = "Polygon"
	p.Perimeter.Coordinates = make([][][]float64, 1)
	p.Perimeter.Coordinates[0] = make([][]float64, 0)

	if len(rq.Perimeter) != 0 {
		for _, point := range rq.Perimeter {
			p.Perimeter.Coordinates[0] = append(p.Perimeter.Coordinates[0], []float64{point.Longitude, point.Latitude})
		}

		p.Perimeter.Coordinates[0] = append(p.Perimeter.Coordinates[0], []float64{rq.Perimeter[0].Longitude, rq.Perimeter[0].Latitude})
	}

	if err := v.Store.Create(ctx, p); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, p)
}

// Recreate creates project docker and stores their information.
// This function is mapped to the path GET /projects/{project_id}/recreate
func (v Projects) Recreate(c echo.Context) error {
	ctx := c.Request().Context()

	projectID := c.Param("project_id")

	_, err := v.Store.Get(ctx, projectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	addr := fmt.Sprintf(
		"amqp://%s:%s@%s:%d/",
		v.Config.Rabbitmq.User,
		v.Config.Rabbitmq.Pass,
		v.Config.Rabbitmq.Host,
		v.Config.Rabbitmq.Port,
	)

	// predefined environment variables
	// This newly created project is under supervision of platform admin
	envs := []runner.Env{
		{Name: "DB_URL", Value: v.Config.Database.URL},
		{Name: "BROKER_URL", Value: addr},
		{Name: "OWNER", Value: "platform.avidnetco@gmail.com"},
	}

	// let's create new dockers for the project
	r, err := v.Manager.New(ctx, projectID, envs)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	p, err := v.Store.Update(ctx, projectID, map[string]interface{}{
		"runner": r,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, p)
}

// Show gets the data for one project. This function is mapped to
// the path GET /projects/{project_id}
func (v Projects) Show(c echo.Context) error {
	ctx := c.Request().Context()

	projectID := c.Param("project_id")

	p, err := v.Store.Get(ctx, projectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	ins, err := v.Manager.Show(ctx, p.Runner)
	if err != nil {
		// There is no available docker for this project. We do not return an error in this condition,
		// but the user must find out based on empty inspect that he must recreate the docker for this project,
		// or try to find out better details from the Portainer.
	} else {
		p.Inspects = ins
	}

	return c.JSON(http.StatusOK, p)
}

// Update updates the name of the project. The project owner is passed as an environment variable
// to project docker so it cannot be changed.
// This function is mapped to the path PUT /projects/{project_id}
func (v Projects) Update(c echo.Context) error {
	ctx := c.Request().Context()

	projectID := c.Param("project_id")

	var name string
	if err := c.Bind(&name); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	p, err := v.Store.Update(ctx, projectID, map[string]interface{}{
		"name": name,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, p)
}

// Destroy deletes a project from the DB and its docker. This function is mapped
// to the path DELETE /projects/{project_id}
func (v Projects) Destroy(c echo.Context) error {
	ctx := c.Request().Context()

	projectID := c.Param("project_id")

	p, err := v.Store.Get(ctx, projectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// remove project runner
	if err := v.Manager.Remove(ctx, p.Runner); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := v.Store.Delete(ctx, projectID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, p)
}

// Logs returns project execution logs and errors. This function is mapped
// to the path GET /projects/{project_id}/logs
func (v Projects) Logs(c echo.Context) error {
	ctx := c.Request().Context()

	projectID := c.Param("project_id")

	limit, err := strconv.Atoi(c.Param("limit"))
	if err != nil {
		limit = 10 // default limit is 10
	}

	pls, err := v.Store.Logs(ctx, projectID, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, pls)
}
