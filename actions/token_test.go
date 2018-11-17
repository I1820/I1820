/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 17-11-2018
 * |
 * | File Name:     token_test.go
 * +===============================================
 */

package actions

import (
	"github.com/I1820/pm/models"
	"github.com/I1820/types"
)

func (as *ActionSuite) Test_TokensResource_Create_Destroy() {
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

	// Create (POST /api/projects/{project_id}/things/{thing_id}/tokens)
	var tk types.Thing
	resk := as.JSON("/api/projects/%s/things/%s/tokens", pID, tID).Get()
	as.Equalf(200, resk.Code, "Error: %s", resk.Body.String())
	resk.Bind(&tk)
	as.Equal(len(tk.Tokens), 2)

	keyOriginal := tk.Tokens[0]
	keyNew := tk.Tokens[1]

	// Destroy (DELETE /api/projects/{project_id}/things/{thing_id}/tokens/{token})
	resf := as.JSON("/api/projects/%s/things/%s/tokens/%s", pID, tID, keyNew).Get()
	as.Equalf(200, resf.Code, "Error: %s", resf.Body.String())
	resf.Bind(&tk)
	as.Equal(len(tk.Tokens), 1)
	as.Equal(tk.Tokens[0], keyOriginal)

	// Destroy thing
	resr := as.JSON("/api/projects/%s/things/%s", pID, tID).Delete()
	as.Equalf(200, resr.Code, "Error: %s", resr.Body.String())

	// Destroy project
	resd := as.JSON("/api/projects/%s", pID).Delete()
	as.Equalf(200, resd.Code, "Error: %s", resd.Body.String())
}
