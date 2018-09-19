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

var tID = ""

func (as *ActionSuite) Test_ThingsResource_Create() {
	// Create project
	var pr models.Project
	resc := as.JSON("/api/projects").Post(projectReq{Name: pName, Owner: pOwner})
	as.Equalf(200, resc.Code, "Error: %s", resc.Body.String())
	resc.Bind(&pr)
	pID = pr.ID

	// Create thing (POST /api/projects/{project_id}/things)
	var tc models.Thing
	rest := as.JSON("/api/projects/%s/things", pID).Post(thingReq{Name: tName})
	as.Equalf(200, rest.Code, "Error: %s", rest.Body.String())
	rest.Bind(&tc)
	tID = tc.ID

	// Show (GET /api/projects/{project_id}/things/{thing_id}
	var ts models.Thing
	ress := as.JSON("/api/projects/%s/things/%s", pID, tID).Get()
	as.Equalf(200, ress.Code, "Error: %s", ress.Body.String())
	ress.Bind(&ts)

	as.Equal(ts, tc)

	// List (GET /api/projects/{project_id}/things)
	resl := as.JSON("/api/projects/%s/things", pID).Get()
	as.Equalf(200, resl.Code, "Error: %s", resl.Body.String())

	// Destroy project
	resd := as.JSON("/api/projects/%s", pID).Delete()
	as.Equalf(200, resd.Code, "Error: %s", resd.Body.String())
}
