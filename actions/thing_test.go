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

import (
	"github.com/I1820/pm/models"
	"github.com/I1820/types"
)

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
	var tc types.Thing
	// build thing creation request
	var rq thingReq
	rq.Name = tName
	rq.Location.Latitude = 35.807657 // I1820 location in velenjak
	rq.Location.Longitude = 51.398408
	rest := as.JSON("/api/projects/%s/things", pID).Post(rq)
	as.Equalf(200, rest.Code, "Error: %s", rest.Body.String())
	rest.Bind(&tc)
	tID = tc.ID

	// Show (GET /api/projects/{project_id}/things/{thing_id}
	var ts types.Thing
	ress := as.JSON("/api/projects/%s/things/%s", pID, tID).Get()
	as.Equalf(200, ress.Code, "Error: %s", ress.Body.String())
	ress.Bind(&ts)

	as.Equal(ts, tc)

	// List (GET /api/projects/{project_id}/things)
	resl := as.JSON("/api/projects/%s/things", pID).Get()
	as.Equalf(200, resl.Code, "Error: %s", resl.Body.String())

	// GeoWithin (POST /api/projects/{project_id}/things/geo)
	var tg []types.Thing
	resg := as.JSON("/api/projects/%s/things/geo", pID).Post(geoWithinReq{
		[][]float64{
			[]float64{35.806731, 51.398618},
			[]float64{35.807784, 51.397810},
			[]float64{35.807827, 51.399516},
			[]float64{35.806731, 51.398618},
		},
	})
	as.Equalf(200, resg.Code, "Error: %s", resg.Body.String())
	resg.Bind(&tg)

	as.Equal(len(tg), 1) // el-thing must be found
	as.Equal(tg[0], tc)

	// Destroy (DELETE /api/projects/{project_id}/things)
	resr := as.JSON("/api/projects/%s/things/%s", pID, tID).Delete()
	as.Equalf(200, resr.Code, "Error: %s", resr.Body.String())

	// Destroy project
	resd := as.JSON("/api/projects/%s", pID).Delete()
	as.Equalf(200, resd.Code, "Error: %s", resd.Body.String())
}
