/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 28-04-2018
 * |
 * | File Name:     query_test.go
 * +===============================================
 */

package actions

const thingID = "5ba3f1a287a142b0a840fae1"
const projectID = "5ba3f19c87a142b0a840fae0"

func (as *ActionSuite) Test_QueriesResource_List() {
	var results []listResp

	res := as.JSON("/api/projects/%s/things/%s/queries/list", projectID, thingID).Get()
	as.Equal(200, res.Code)

	res.Bind(&results)

	as.NotEqual(len(results), 0)

	for _, r := range results {
		if r.ID == "100" {
			as.Equal(r.Total, 4)
		}
	}
}
