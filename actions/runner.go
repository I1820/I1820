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
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/aiotrc/pm/models"
	"github.com/gobuffalo/buffalo"
	"github.com/mongodb/mongo-go-driver/bson"
	mgo "github.com/mongodb/mongo-go-driver/mongo"
)

// RunnersHandler sends request to specific GoRunner
func RunnersHandler(c buffalo.Context) error {
	name := c.Param("project_id")
	path := c.Param("path")

	var p models.Project

	dr := db.Collection("pm").FindOne(context.Background(), bson.NewDocument(
		bson.EC.String("name", name),
	))

	if err := dr.Decode(&p); err != nil {
		if err == mgo.ErrNoDocuments {
			return c.Error(http.StatusNotFound, fmt.Errorf("Project %s not found", name))
		}
		return c.Error(http.StatusInternalServerError, err)
	}

	url, err := url.Parse(fmt.Sprintf("http://172.17.0.1:%s/api", p.Runner.Port))
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	c.Request().URL.Path = path
	return buffalo.WrapHandler(
		httputil.NewSingleHostReverseProxy(url),
	)(c)
}
