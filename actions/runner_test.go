/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 17-07-2018
 * |
 * | File Name:     runner_test.go
 * +===============================================
 */

package actions

import (
	"fmt"
	"time"

	"github.com/I1820/pm/models"
)

func (as *ActionSuite) Test_RunnersHandler() {
	// Create project
	var pr models.Project
	resc := as.JSON("/api/projects").Post(projectReq{Name: pName, Owner: pOwner})
	as.Equalf(200, resc.Code, "Error: %s", resc.Body.String())
	resc.Bind(&pr)
	pID = pr.ID

	// wait for ElRunner make ready
	time.Sleep(15 * time.Second)

	// ElRunner About API (GET /api/runners/{project_id}/about)
	res := as.JSON("/api/runners/%s/about", pID).Get()
	fmt.Println(res.Header())
	as.Equalf(200, res.Code, "Error: %s", res.Body.String())
	as.Contains(res.Body.String(), "18.20 is leaving us")

	// Destroy project
	resd := as.JSON("/api/projects/%s", pID).Delete()
	as.Equalf(200, resd.Code, "Error: %s", resd.Body.String())
}
