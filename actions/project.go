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
	"github.com/mongodb/mongo-go-driver/mongo/aggregateopt"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"github.com/mongodb/mongo-go-driver/mongo/mongoopt"

	mgo "github.com/mongodb/mongo-go-driver/mongo"
)

// ProjectsResource manages existing projects
type ProjectsResource struct {
	buffalo.Resource
}

// project request payload
type projectReq struct {
	Name        string            `json:"name" validate:"required"`        // project name
	Owner       string            `json:"owner" validate:"required,email"` // project owner email address
	Envs        map[string]string `json:"envs"`                            // project environment variables
	Description string            `json:"description"`                     // project description
	Perimeter   []struct {        // project operational perimeter
		Latitude  float64 `json:"lat"`
		Longitude float64 `json:"long"`
	} `json:"perimeter"`
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

	if _, err := db.Collection("projects").InsertOne(c, p); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, r.JSON(p))
}

// Recreate creates project docker and stores their information.
// This function is mapped to the path GET /projects/{project_id}/recreate
func (v ProjectsResource) Recreate(c buffalo.Context) error {
	projectID := c.Param("project_id")

	if err := validate.Var(projectID, "required"); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	var p models.Project

	dr := db.Collection("projects").FindOne(c, bson.NewDocument(
		bson.EC.String("name", projectID),
	))

	if err := dr.Decode(&p); err != nil {
		if err == mgo.ErrNoDocuments {
			return c.Error(http.StatusNotFound, fmt.Errorf("Project %s not found", projectID))
		}
		return c.Error(http.StatusInternalServerError, err)
	}

	// predefined environment variables
	// This newely created project is under supervision of platform admin
	envs := []runner.Env{
		{Name: "DB_URL", Value: envy.Get("DB_URL", "mongodb://192.168.72.1:27017")},
		{Name: "BROKER_URL", Value: envy.Get("BROKER_URL", "tcp://192.168.72.1:1883")},
		{Name: "OWNER", Value: "platform.avidnetco@gmail.com"},
	}

	// let's create new dockers for the project
	if err := models.ReProject(c, envs, &p); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	du := db.Collection("projects").FindOneAndUpdate(c, bson.NewDocument(
		bson.EC.String("name", projectID),
	), bson.NewDocument(
		bson.EC.SubDocumentFromElements("$set", bson.EC.Interface("runner", p.Runner)),
	), findopt.ReturnDocument(mongoopt.After))

	if err := du.Decode(&p); err != nil {
		if err == mgo.ErrNoDocuments {
			return c.Error(http.StatusNotFound, fmt.Errorf("Project %s not found", projectID))
		}
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
		bson.EC.String("name", projectID),
	))

	if err := dr.Decode(&p); err != nil {
		if err == mgo.ErrNoDocuments {
			return c.Error(http.StatusNotFound, fmt.Errorf("Project %s not found", projectID))
		}
		return c.Error(http.StatusInternalServerError, err)
	}

	ins, err := p.Runner.Show(c)
	if err != nil {
		// There is no available docker for this project. We do not return an error in this condition,
		// but the user must find out based on empty inspect that he must recreate the docker for this project,
		// or try to find out better details from the Portainer.
	} else {
		p.Inspects = ins
	}

	return c.Render(http.StatusOK, r.JSON(p))
}

// Update updates the name of the project. The project owner is passed as an environment variable
// to project docker so it cannot be changed.
// This function is mapped to the path PUT /projects/{project_id}
func (v ProjectsResource) Update(c buffalo.Context) error {
	projectID := c.Param("project_id")

	var name string
	if err := c.Bind(&name); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	if err := validate.Var(name, "required"); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	var p models.Project

	dr := db.Collection("projects").FindOneAndUpdate(c, bson.NewDocument(
		bson.EC.String("name", projectID),
	), bson.NewDocument(
		bson.EC.SubDocumentFromElements("$set", bson.EC.String("name", name)),
	), findopt.ReturnDocument(mongoopt.After))

	if err := dr.Decode(&p); err != nil {
		if err == mgo.ErrNoDocuments {
			return c.Error(http.StatusNotFound, fmt.Errorf("Project %s not found", projectID))
		}
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, r.JSON(p))
}

// Destroy deletes a project from the DB and its docker. This function is mapped
// to the path DELETE /projects/{project_id}
func (v ProjectsResource) Destroy(c buffalo.Context) error {
	projectID := c.Param("project_id")

	var p models.Project

	dr := db.Collection("projects").FindOne(c, bson.NewDocument(
		bson.EC.String("name", projectID),
	))

	if err := dr.Decode(&p); err != nil {
		if err == mgo.ErrNoDocuments {
			return c.Error(http.StatusNotFound, fmt.Errorf("Project %s not found", projectID))
		}
		return c.Error(http.StatusInternalServerError, err)
	}

	// remove project runner
	if err := p.Runner.Remove(c); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	// remove project entity from database
	if _, err := db.Collection("projects").DeleteOne(c, bson.NewDocument(
		bson.EC.String("name", projectID),
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
		limit = 10 // default limit is 10
	}

	cur, err := db.Collection(fmt.Sprintf("projects.logs.%s", projectID)).Aggregate(c, bson.NewArray(
		bson.VC.DocumentFromElements(
			bson.EC.SubDocumentFromElements("$sort", bson.EC.Int32("Time", -1)),
		),
		bson.VC.DocumentFromElements(
			bson.EC.Int32("$limit", int32(limit)),
		),
	), aggregateopt.AllowDiskUse(true))
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
