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

	"github.com/aiotrc/pm/models"
	"github.com/mongodb/mongo-go-driver/bson"
)

const pName = "Her"

func (as *ActionSuite) Test_ProjectsResource_Create() {
	var pr models.Project

	res := as.JSON("/api/projects").Post(projectReq{Name: pName})
	as.Equalf(200, res.Code, "Error: %s", res.Body.String())
	res.Bind(&pr)

	var pd models.Project
	dr := db.Collection("pm").FindOne(context.Background(), bson.NewDocument(
		bson.EC.String("name", pName),
	))

	as.NoError(dr.Decode(&pd))

	as.Equal(pd, pr)
}

func (as *ActionSuite) Test_ProjectsResource_Show() {
	var p models.Project

	res := as.JSON("/api/projects/%s", pName).Get()
	as.Equalf(200, res.Code, "Error: %s", res.Body.String())
	res.Bind(&p)
}
