/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 18-10-2018
 * |
 * | File Name:     tag_test.go
 * +===============================================
 */

package actions

import (
	"github.com/I1820/pm/models"
	"github.com/I1820/types"
)

const tag1 = "Elahe"
const tag2 = "Leaving"
const tag3 = "Us"

func (as *ActionSuite) Test_TagsResource_Create() {
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

	// Create (POST /api/projects/{project_id}/things/{thing_id}/tags)
	resg := as.JSON("/api/projects/%s/things/%s/tags", pID, tID).Post(tagReq{Tags: []string{tag1, tag2, tag3}})
	as.Equalf(200, resg.Code, "Error: %s", resg.Body.String())

	// List (Get /api/projects/{project_id}/things/{thing_id}/tags)
	var tags []string
	resl := as.JSON("/api/projects/%s/things/%s/tags", pID, tID).Get()
	as.Equalf(200, resl.Code, "Error: %s", resl.Body.String())
	resl.Bind(&tags)
	as.Equal(tags[0], tag1)
	as.Equal(tags[1], tag2)
	as.Equal(tags[2], tag3)

	// Destroy thing
	resr := as.JSON("/api/projects/%s/things/%s", pID, tID).Delete()
	as.Equalf(200, resr.Code, "Error: %s", resr.Body.String())

	// Destroy project
	resd := as.JSON("/api/projects/%s", pID).Delete()
	as.Equalf(200, resd.Code, "Error: %s", resd.Body.String())
}
