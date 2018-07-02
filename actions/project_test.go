/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 02-07-2018
 * |
 * | File Name:     actions/project_test.go
 * +===============================================
 */

package actions

import "github.com/aiotrc/pm/project"

func (as *ActionSuite) Test_ProjectNewHandler() {
	var p *project.Project

	res := as.JSON("/api/project").Post(projectReq{Name: "Her"})
	as.Equal(200, res.Code)
	res.Bind(p)

	as.Info(p)
}
