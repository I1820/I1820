/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 02-07-2018
 * |
 * | File Name:     actions/project_test.go
 * +===============================================
 */

package actions

import (
	"context"

	"github.com/I1820/pm/models"
	"github.com/mongodb/mongo-go-driver/bson"
)

const pName = "kj"
const pOwner = "parham.alvani@gmail.com"

var pID = ""

func (as *ActionSuite) Test_ProjectsResource_Create_Show_Destroy() {
	var pr models.Project

	// Create (POST /api/projects)
	resc := as.JSON("/api/projects").Post(projectReq{Name: pName, Owner: pOwner})
	as.Equalf(200, resc.Code, "Error: %s", resc.Body.String())
	resc.Bind(&pr)
	pID = pr.ID

	// check database for project existence
	var pd models.Project
	dr := db.Collection("projects").FindOne(context.Background(), bson.NewDocument(
		bson.EC.String("id", pID),
	))

	as.NoError(dr.Decode(&pd))

	as.Equal(pd, pr)

	// Show (GET /api/projects/{project_id})
	ress := as.JSON("/api/projects/%s", pID).Get()
	as.Equalf(200, ress.Code, "Error: %s", ress.Body.String())
	ress.Bind(&pr)
	pr.Inspects = nil

	as.Equal(pd, pr)

	// Destroy (DELETE /api/projects/{project_id})
	resd := as.JSON("/api/projects/%s", pID).Delete()
	as.Equalf(200, resd.Code, "Error: %s", resd.Body.String())
	resd.Bind(&pr)

	as.Equal(pd, pr)
}

func (as *ActionSuite) Test_ProjectsResource_List() {
	var ps []models.Project

	res := as.JSON("/api/projects").Get()
	as.Equalf(200, res.Code, "Error: %s", res.Body.String())
	res.Bind(&ps)
}
