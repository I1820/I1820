package handler

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/I1820/I1820/runner"
	"github.com/I1820/I1820/store"
	"github.com/labstack/echo/v4"
)

type Runner struct {
	Store   store.Project
	Manager *runner.Manager

	DockerHost string
}

func (r Runner) Register(g *echo.Group) {
	g.GET("/runners/pull", r.Pull)
	g.Any("/runners/{project_id}/{path:.+}", r.PassThrough)
}

// PassThrough sends request to specific Runner
func (r Runner) PassThrough(c echo.Context) error {
	ctx := c.Request().Context()

	projectID := c.Param("project_id")
	path := c.Param("path")

	p, err := r.Store.Get(ctx, projectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	url, err := url.Parse(
		fmt.Sprintf("http://%s:%s/", r.DockerHost, p.Runner.Port),
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	c.Request().URL.Path = strings.TrimSuffix(path, "/")

	return echo.WrapHandler(
		httputil.NewSingleHostReverseProxy(url),
	)(c)
}

// Pull pulls the latest version of required images.
func (r Runner) Pull(c echo.Context) error {
	ctx := c.Request().Context()

	rs, err := r.Manager.Pull(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, rs)
}
