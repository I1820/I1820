/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 17-07-2018
 * |
 * | File Name:     runner.go
 * +===============================================
 */

package actions

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/I1820/pm/models"
	"github.com/I1820/pm/runner"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/mongodb/mongo-go-driver/bson"
	mgo "github.com/mongodb/mongo-go-driver/mongo"
)

// RunnersHandler sends request to specific ElRunner
func RunnersHandler(c buffalo.Context) error {
	projectID := c.Param("project_id")
	path := c.Param("path")

	var p models.Project

	dr := db.Collection("projects").FindOne(c, bson.NewDocument(
		bson.EC.String("_id", projectID),
	))

	if err := dr.Decode(&p); err != nil {
		if err == mgo.ErrNoDocuments {
			return c.Error(http.StatusNotFound, fmt.Errorf("Project %s not found", projectID))
		}
		return c.Error(http.StatusInternalServerError, err)
	}

	url, err := url.Parse(fmt.Sprintf("http://%s:%s/", envy.Get("D_HOST", "127.0.0.1"), p.Runner.Port))
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	c.Request().URL.Path = strings.TrimSuffix(path, "/")
	return buffalo.WrapHandler(
		httputil.NewSingleHostReverseProxy(url),
	)(c)
}

// PullHandler pulls the latest version of required images.
func PullHandler(c buffalo.Context) error {
	rs, err := runner.Pull(c)
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}
	return c.Render(http.StatusOK, r.JSON(rs))
}
