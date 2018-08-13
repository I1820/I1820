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

func (as *ActionSuite) Test_QueriesResource_List() {
	var results []listResp

	res := as.JSON("/api/queries/list").Get()
	as.Equal(200, res.Code)

	res.Bind(&results)

	as.Equal(len(results), 0)

	for _, r := range results {
		if r.ID == "0000000000000003" {
			as.Equal(r.Total, 1)
		}
	}
}
