/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 07-07-2018
 * |
 * | File Name:     thing_test.go
 * +===============================================
 */

package actions

import "github.com/I1820/pm/models"

const tName = "0000000000000073"

func (as *ActionSuite) Test_ThingsResource_Create() {
	// Create
	resc := as.JSON("/api/%s/projects", uName).Post(projectReq{Name: pName})
	as.Equalf(200, resc.Code, "Error: %s", resc.Body.String())

	// Create
	var tc models.Thing
	rest := as.JSON("/api/things").Post(thingReq{Name: tName, Project: pName, User: uName})
	as.Equalf(200, rest.Code, "Error: %s", rest.Body.String())
	rest.Bind(&tc)

	// Show
	var ts models.Thing
	ress := as.JSON("/api/things/%s", tName).Get()
	as.Equalf(200, ress.Code, "Error: %s", ress.Body.String())
	ress.Bind(&ts)

	as.Equal(ts, tc)

	// Destroy
	resd := as.JSON("/api/%s/projects/%s", uName, pName).Delete()
	as.Equalf(200, resd.Code, "Error: %s", resd.Body.String())
}

func (as *ActionSuite) Test_ThingsResource_List() {
	res := as.JSON("/api/things").Get()
	as.Equalf(200, res.Code, "Error: %s", res.Body.String())

}
