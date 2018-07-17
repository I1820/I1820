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

func (as *ActionSuite) Test_RunnersHandler() {
	// Create
	resc := as.JSON("/api/projects").Post(projectReq{Name: pName})
	as.Equalf(200, resc.Code, "Error: %s", resc.Body.String())

	// GoRunner About
	res := as.JSON("/api/runners/%s/about", pName).Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "18.20 is leaving us")

	// Destroy
	/*
		resd := as.JSON("/api/projects/%s", pName).Delete()
		as.Equalf(200, resd.Code, "Error: %s", resd.Body.String())
	*/
}
