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

import "time"

func (as *ActionSuite) Test_RunnersHandler() {
	// Create
	resc := as.JSON("/api/%s/projects", uName).Post(projectReq{Name: pName})
	as.Equalf(200, resc.Code, "Error: %s", resc.Body.String())

	// wait for GoRunner
	time.Sleep(15 * time.Second)

	// GoRunner About
	res := as.JSON("/api/%s/runners/%s/about", uName, pName).Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "18.20 is leaving us")

	// Destroy
	resd := as.JSON("/api/%s/projects/%s", uName, pName).Delete()
	as.Equalf(200, resd.Code, "Error: %s", resd.Body.String())
}
