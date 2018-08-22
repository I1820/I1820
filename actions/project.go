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

// ProjectsResource manages existing projects
type ProjectsResource struct {
	buffalo.Resource
}

// project request payload
type projectReq struct {
	Name string            `json:"name"`
	Envs map[string]string `json:"envs"`
	// TODO adds docker constraints
}

// List gets all projects. This function is mapped to the path
// GET /projects
func (v ProjectsResource) List(c buffalo.Context) error {
	ps := make([]models.Project, 0)

	cur, err := db.Collection("projects").Find(c, bson.NewDocument())
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
// path POST /projects
func (v ProjectsResource) Create(c buffalo.Context) error {
	var rq projectReq
	if err := c.Bind(&rq); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	if rq.Name == "" {
		return c.Error(http.StatusBadRequest, fmt.Errorf("Name should not be null"))
	}
	name := rq.Name

	envs := []runner.Env{
		{Name: "MONGO_URL", Value: envy.Get("DB_URL", "mongodb://172.18.0.1:27017")},
	}

	for envKey, envVal := range rq.Envs {
		envs = append(envs, runner.Env{Name: envKey, Value: envVal})
	}

	p, err := models.NewProject(c, name, envs)
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	if _, err := db.Collection("projects").InsertOne(c, p); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	// numberOfCreatedProjects.Inc()

	return c.Render(http.StatusOK, r.JSON(p))
}

// Show gets the data for one project. This function is mapped to
// the path GET /projects/{project_id}
func (v ProjectsResource) Show(c buffalo.Context) error {
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

	ins, err := p.Runner.Show(c)
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}
	p.Inspects = ins

	return c.Render(http.StatusOK, r.JSON(p))
}

// Destroy deletes a project from the DB and its docker. This function is mapped
// to the path DELETE /projects/{project_id}
func (v ProjectsResource) Destroy(c buffalo.Context) error {
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
	)); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, r.JSON(p))
}

// Activation activates/deactivates project. This function is mapped
// to the path GET /projects/{project_id}/{t:(?:activate|deactivate)}
func (v ProjectsResource) Activation(c buffalo.Context) error {
	name := c.Param("project_id")

	t := c.Param("t")
	status := false
	if t == "activate" {
		status = true
	}

	dr := db.Collection("projects").FindOneAndUpdate(c, bson.NewDocument(
		bson.EC.String("name", name),
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
// to the path GET /projects/{project_id}/logs
func (v ProjectsResource) Logs(c buffalo.Context) error {
	var pls = make([]models.ProjectLog, 0)

	id := c.Param("project_id")

	limit, err := strconv.Atoi(c.Param("limit"))
	if err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	cur, err := db.Collection("projects.logs").Aggregate(c, bson.NewArray(
		bson.VC.DocumentFromElements(
			bson.EC.SubDocumentFromElements("$match", bson.EC.String("project", id)),
		),
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
