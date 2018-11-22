/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 24-09-2018
 * |
 * | File Name:     conn_test.go
 * +===============================================
 */

package actions

import (
	"github.com/I1820/pm/models"
	"github.com/I1820/types"
	"github.com/I1820/types/connectivity"
)

const cName = "ttn"

func (as *ActionSuite) Test_ConnectivitiesResource_Create() {
	// Create project
	var pr models.Project
	resc := as.JSON("/api/projects").Post(projectReq{Name: pName, Owner: pOwner})
	as.Equalf(200, resc.Code, "Error: %s", resc.Body.String())
	resc.Bind(&pr)
	pID = pr.ID

	// Create thing
	var tc types.Thing
	rest := as.JSON("/api/projects/%s/things", pID).Post(thingReq{Name: tName})
	as.Equalf(200, rest.Code, "Error: %s", rest.Body.String())
	rest.Bind(&tc)
	tID = tc.ID

	// Create (POST /api/projects/{project_id}/things/{thing_id}/connectivities)
	resa := as.JSON("/api/projects/%s/things/%s/connectivities", tID, pID).Post(connectivityReq{Name: cName, Info: connectivity.TTN{ApplicationID: "fan", DeviceEUI: "000AE31955C049FC"}})
	as.Equalf(200, resa.Code, "Error: %s", resa.Body.String())

	// Show (Get /api/projects/{project_id}/things/{thing_id}/connectivities/{connectivity_name})
	var ttn connectivity.TTN
	ress := as.JSON("/api/projects/%s/things/%s/connectivities/%s", tID, pID, cName).Get()
	as.Equalf(200, ress.Code, "Error: %s", ress.Body.String())
	ress.Bind(&ttn)
	as.Equal("fan", ttn.ApplicationID)
	as.Equal("000AE31955C049FC", ttn.DeviceEUI)

	// Destroy thing
	resr := as.JSON("/api/projects/%s/things/%s", pID, tID).Delete()
	as.Equalf(200, resr.Code, "Error: %s", resr.Body.String())

	// Destroy project
	resd := as.JSON("/api/projects/%s", pID).Delete()
	as.Equalf(200, resd.Code, "Error: %s", resd.Body.String())
}
