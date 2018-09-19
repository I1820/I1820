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
	"github.com/mongodb/mongo-go-driver/bson/objectid"

	mgo "github.com/mongodb/mongo-go-driver/mongo"
)

// ProjectsResource manages existing projects
type ProjectsResource struct {
	buffalo.Resource
}

// project request payload
type projectReq struct {
	Name  string            `json:"name" validate:"required"`        // project name
	Owner string            `json:"owner" validate:"required,email"` // project owner email address
	Envs  map[string]string `json:"envs"`                            // project environment variables
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

	if err := validate.Struct(rq); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	// predefined environment variables
	envs := []runner.Env{
		{Name: "DB_URL", Value: envy.Get("DB_URL", "mongodb://192.168.72.1:27017")},
		{Name: "BROKER_URL", Value: envy.Get("BROKER_URL", "tcp://192.168.72.1:1883")},
		{Name: "OWNER", Value: rq.Owner},
	}

	// user-defined environment variables
	for envKey, envVal := range rq.Envs {
		envs = append(envs, runner.Env{Name: envKey, Value: envVal})
	}

	id := objectid.New().Hex()

	// creates project entity with its docker (have fun :D)
	p, err := models.NewProject(c, id, rq.Name, envs)
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}
	p.ID = id

	if _, err := db.Collection("projects").InsertOne(c, p); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, r.JSON(p))
}

// Show gets the data for one project. This function is mapped to
// the path GET /projects/{project_id}
func (v ProjectsResource) Show(c buffalo.Context) error {
	projectID := c.Param("project_id")

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
	projectID := c.Param("project_id")

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

	if err := p.Runner.Remove(c); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	if _, err := db.Collection("projects").DeleteOne(c, bson.NewDocument(
		bson.EC.String("_id", projectID),
	)); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	if err := db.Collection(fmt.Sprintf("projects.logs.%s", projectID)).Drop(c); err != nil {
	}

	return c.Render(http.StatusOK, r.JSON(p))
}

// Logs returns project execution logs and errors. This function is mapped
// to the path GET /projects/{project_id}/logs
func (v ProjectsResource) Logs(c buffalo.Context) error {
	var pls = make([]models.ProjectLog, 0)

	projectID := c.Param("project_id")

	limit, err := strconv.Atoi(c.Param("limit"))
	if err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	cur, err := db.Collection(fmt.Sprintf("projects.logs.%s", projectID)).Aggregate(c, bson.NewArray(
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
