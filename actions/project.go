/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 02-07-2018
 * |
 * | File Name:     actions/project.go
 * +===============================================
 */

package actions

import (
	"net/http"

	"github.com/aiotrc/pm/project"
	"github.com/aiotrc/pm/runner"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
)

// project request payload
type projectReq struct {
	Name string `json:"name" binding:"required"`
	// TODO adds docker constraints and envs
}

func projectNewHandler(c buffalo.Context) error {
	var rq projectReq
	if err := c.Bind(&rq); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	name := rq.Name

	p, err := project.New(name, []runner.Env{
		{Name: "MONGO_URL", Value: envy.Get("DB_URL", "mongodb://172.18.0.1:27017")},
	})
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	if _, err := db.Collection("pm").InsertOne(c, p); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	// numberOfCreatedProjects.Inc()

	return c.Render(http.StatusOK, r.JSON(p))
}
