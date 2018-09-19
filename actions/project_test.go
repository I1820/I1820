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

const pName = "Her"
const uName = "1995parham"

func (as *ActionSuite) Test_ProjectsResource_Create_Show_Destroy() {
	var pr models.Project

	// Create
	resc := as.JSON("/api/%s/projects", uName).Post(projectReq{Name: pName})
	as.Equalf(200, resc.Code, "Error: %s", resc.Body.String())
	resc.Bind(&pr)

	// Database
	var pd models.Project
	dr := db.Collection("projects").FindOne(context.Background(), bson.NewDocument(
		bson.EC.String("name", pName),
		bson.EC.String("user", uName),
	))

	as.NoError(dr.Decode(&pd))

	as.Equal(pd, pr)

	// Show
	ress := as.JSON("/api/%s/projects/%s", uName, pName).Get()
	as.Equalf(200, ress.Code, "Error: %s", ress.Body.String())
	ress.Bind(&pr)
	pr.Inspects = nil

	as.Equal(pd, pr)

	// Deactivate
	resa := as.JSON("/api/%s/projects/%s/deactivate", uName, pName).Get()
	as.Equalf(200, resa.Code, "Error: %s", resa.Body.String())
	resa.Bind(&pd)

	// Destroy
	resd := as.JSON("/api/%s/projects/%s", uName, pName).Delete()
	as.Equalf(200, resd.Code, "Error: %s", resd.Body.String())
	resd.Bind(&pr)

	as.Equal(pd, pr)
}

func (as *ActionSuite) Test_ProjectsResource_List() {
	var ps []models.Project

	res := as.JSON("/api/%s/projects", uName).Get()
	as.Equalf(200, res.Code, "Error: %s", res.Body.String())
	res.Bind(&ps)
}
