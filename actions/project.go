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
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/I1820/pm/models"
	"github.com/I1820/pm/runner"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/mongodb/mongo-go-driver/bson"
	mgo "github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"github.com/mongodb/mongo-go-driver/mongo/mongoopt"
)

var nameRegexp *regexp.Regexp

func init() {
	rg, err := regexp.Compile("[0-9a-zA-Z]")
	if err == nil {
		nameRegexp = rg
	}
}

// UserID acts as a middleware and validates userid
func UserID(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		user := c.Param("user_id")
		if !nameRegexp.MatchString(user) {
			return c.Error(http.StatusBadRequest, fmt.Errorf("Invalid user id: %s", user))
		}
		return next(c)
	}
}

// ProjectsResource manages existing projects
type ProjectsResource struct {
	buffalo.Resource
}

// project request payload
type projectReq struct {
	Name string            `json:"name"` // project_id
	Envs map[string]string `json:"envs"` // project environment variables

	// TODO adds docker constraints
}

// List gets all projects. This function is mapped to the path
// GET /{user_id}/projects
func (v ProjectsResource) List(c buffalo.Context) error {
	user := c.Param("user_id")
	ps := make([]models.Project, 0)

	cur, err := db.Collection("projects").Find(c, bson.NewDocument(
		bson.EC.String("user", user),
	))
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	for cur.Next(c) {
		var p models.Project

		if err := cur.Decode(&p); err != nil {
			return c.Error(http.StatusInternalServerError, err)
		}

		ps = append(ps, p)
	}
	if err := cur.Close(c); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, r.JSON(ps))
}

// Create adds a project to the DB and creates its docker. This function is mapped to the
// path POST /{user_id}/projects
func (v ProjectsResource) Create(c buffalo.Context) error {
	user := c.Param("user_id")

	var rq projectReq
	if err := c.Bind(&rq); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	if rq.Name == "" || !nameRegexp.MatchString(rq.Name) {
		return c.Error(http.StatusBadRequest, fmt.Errorf("Invalid name: %s", rq.Name))
	}
	name := rq.Name

	// predefined environment variables
	envs := []runner.Env{
		{Name: "DB_URL", Value: envy.Get("DB_URL", "mongodb://192.168.72.1:27017")},
		{Name: "BROKER_URL", Value: envy.Get("BROKER_URL", "tcp://192.168.72.1:1883")},
		{Name: "USER", Value: user},
	}

	// user-defined environment variables
	for envKey, envVal := range rq.Envs {
		envs = append(envs, runner.Env{Name: envKey, Value: envVal})
	}

	// creates project entity with its docker (have fun :D)
	p, err := models.NewProject(c, user, name, envs)
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	if _, err := db.Collection("projects").InsertOne(c, p); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, r.JSON(p))
}

// Show gets the data for one project. This function is mapped to
// the path GET /{user_id}/projects/{project_id}
func (v ProjectsResource) Show(c buffalo.Context) error {
	user := c.Param("user_id")
	name := c.Param("project_id")

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

	ins, err := p.Runner.Show(c)
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}
	p.Inspects = ins

	return c.Render(http.StatusOK, r.JSON(p))
}

// Destroy deletes a project from the DB and its docker. This function is mapped
// to the path DELETE /{user_id}/projects/{project_id}
func (v ProjectsResource) Destroy(c buffalo.Context) error {
	user := c.Param("user_id")
	name := c.Param("project_id")

	var p models.Project

	dr := db.Collection("projects").FindOne(c, bson.NewDocument(
		bson.EC.String("name", name),
	))

	if err := dr.Decode(&p); err != nil {
		if err == mgo.ErrNoDocuments {
			return c.Error(http.StatusNotFound, fmt.Errorf("Project %s not found", name))
		}
		return c.Error(http.StatusInternalServerError, err)
	}

	if err := p.Runner.Remove(c); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	if _, err := db.Collection("projects").DeleteOne(c, bson.NewDocument(
		bson.EC.String("name", name),
		bson.EC.String("user", user),
	)); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	if err := db.Collection(fmt.Sprintf("projects.logs.%s", name)).Drop(c); err != nil {
	}

	return c.Render(http.StatusOK, r.JSON(p))
}

// Activation activates/deactivates project. This function is mapped
// to the path GET /{user_id}/projects/{project_id}/{t:(?:activate|deactivate)}
func (v ProjectsResource) Activation(c buffalo.Context) error {
	name := c.Param("project_id")
	user := c.Param("user_id")

	t := c.Param("t")
	status := false
	if t == "activate" {
		status = true
	}

	dr := db.Collection("projects").FindOneAndUpdate(c, bson.NewDocument(
		bson.EC.String("name", name),
		bson.EC.String("user", user),
	), bson.NewDocument(
		bson.EC.SubDocumentFromElements("$set", bson.EC.Boolean("status", status)),
	), findopt.ReturnDocument(mongoopt.After))

	var p models.Project

	if err := dr.Decode(&p); err != nil {
		if err == mgo.ErrNoDocuments {
			return c.Error(http.StatusNotFound, fmt.Errorf("Project %s not found", name))
		}
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, r.JSON(p))
}

// Logs returns project execution logs and errors. This function is mapped
// to the path GET /{user_id}/projects/{project_id}/logs
func (v ProjectsResource) Logs(c buffalo.Context) error {
	var pls = make([]models.ProjectLog, 0)

	name := c.Param("project_id")
	user := c.Param("user_id")

	limit, err := strconv.Atoi(c.Param("limit"))
	if err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	cur, err := db.Collection(fmt.Sprintf("projects.logs.%s_%s", name, user)).Aggregate(c, bson.NewArray(
		bson.VC.DocumentFromElements(
			bson.EC.Int32("$limit", int32(limit)),
		),
		bson.VC.DocumentFromElements(
			bson.EC.SubDocumentFromElements("$sort", bson.EC.Int32("Time", -1)),
		),
	))
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	for cur.Next(c) {
		var pl models.ProjectLog

		if err := cur.Decode(&pl); err != nil {
			return c.Error(http.StatusInternalServerError, err)
		}

		pls = append(pls, pl)
	}
	if err := cur.Close(c); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, r.JSON(pls))
}
