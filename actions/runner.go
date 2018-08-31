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

	"github.com/I1820/pm/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/mongodb/mongo-go-driver/bson"
	mgo "github.com/mongodb/mongo-go-driver/mongo"
)

// RunnersHandler sends request to specific GoRunner
func RunnersHandler(c buffalo.Context) error {
	name := c.Param("project_id")
	user := c.Param("user_id")
	path := c.Param("path")

	var p models.Project

	dr := db.Collection("projects").FindOne(c, bson.NewDocument(
		bson.EC.String("name", name),
		bson.EC.String("user", user),
	))

	if err := dr.Decode(&p); err != nil {
		if err == mgo.ErrNoDocuments {
			return c.Error(http.StatusNotFound, fmt.Errorf("Project %s not found", name))
		}
		return c.Error(http.StatusInternalServerError, err)
	}

	url, err := url.Parse(fmt.Sprintf("http://%s:%s", envy.Get("D_HOST", "127.0.0.1"), p.Runner.Port))
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	c.Request().URL.Path = path
	return buffalo.WrapHandler(
		httputil.NewSingleHostReverseProxy(url),
	)(c)
}
