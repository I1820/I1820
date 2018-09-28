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

import "time"

const thingID = "5ba3f1a287a142b0a840fae1"
const projectID = "5ba3f19c87a142b0a840fae0"

func (as *ActionSuite) Test_QueriesResource_List() {
	var results []listResp

	res := as.JSON("/api/projects/%s/things/%s/queries/list", projectID, thingID).Get()
	as.Equal(200, res.Code)

	res.Bind(&results)

	as.NotEqual(0, len(results))

	for _, r := range results {
		if r.ID == "100" {
			as.Equal(4, r.Total)
		}
	}
}

func (as *ActionSuite) Test_QueriesResource_PFetch() {
	var results []pfetchResp

	var req fetchReq
	req.Range.To = time.Date(2019, time.September, 11, 0, 0, 0, 0, time.UTC)
	req.Range.From = time.Date(2017, time.September, 11, 0, 0, 0, 0, time.UTC)
	req.Target = "100"
	req.Window.Size = 1

	res := as.JSON("/api/projects/%s/things/%s/queries/pfetch", projectID, thingID).Post(req)
	as.Equal(200, res.Code)

	res.Bind(&results)

	as.Equal(1, len(results))
	as.Equal(6750.0, results[0].Data)
}
