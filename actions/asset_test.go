/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 20-09-2018
 * |
 * | File Name:     asset_test.go
 * +===============================================
 */

package actions

import (
	"github.com/I1820/pm/models"
	"github.com/I1820/types"
)

const aName = "101"
const aTitle = "Temperature"
const aKind = "sensor"
const aType = "number"

func (as *ActionSuite) Test_AssetsResource_Create() {
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

	// Create (POST /api/projects/{project_id}/things/{thing_id}/assets)
	resa := as.JSON("/api/projects/%s/things/%s/assets", pID, tID).Post(assetReq{Name: aName, Title: aTitle, Type: aType, Kind: aKind})
	as.Equalf(200, resa.Code, "Error: %s", resa.Body.String())

	// Show (Get /api/projects/{project_id}/things/{thing_id}/assets/{asset_name})
	var a types.Asset
	ress := as.JSON("/api/projects/%s/things/%s/assets/%s", pID, tID, aName).Get()
	as.Equalf(200, ress.Code, "Error: %s", ress.Body.String())
	ress.Bind(&a)
	as.Equal(aTitle, a.Title)
	as.Equal(aType, a.Type)
	as.Equal(aKind, a.Kind)

	// Destroy thing
	resr := as.JSON("/api/projects/%s/things/%s", pID, tID).Delete()
	as.Equalf(200, resr.Code, "Error: %s", resr.Body.String())

	// Destroy project
	resd := as.JSON("/api/projects/%s", pID).Delete()
	as.Equalf(200, resd.Code, "Error: %s", resd.Body.String())
}
