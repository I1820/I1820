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

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/I1820/types"
	"github.com/labstack/echo/v4"
)

const thingID = "5ba3f1a287a142b0a840fae1"
const projectID = "5ba3f19c87a142b0a840fae0"

func (suite *DMTestSuite) Test_QueriesHandler_List() {
	var results []listResp

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/projects/%s/things/%s/queries/list", projectID, thingID), nil)
	suite.NoError(err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)

	suite.Equal(200, w.Code)

	suite.NoError(json.Unmarshal(w.Body.Bytes(), &results))

	suite.NotEqual(0, len(results))

	for _, r := range results {
		if r.ID == "100" {
			suite.Equal(4, r.Total)
		}
	}
}

func (suite *DMTestSuite) Test_QueriesHandler_PFetch() {
	var results []pfetchResp

	var freq fetchReq
	freq.Range.To = time.Date(2019, time.September, 11, 0, 0, 0, 0, time.UTC)
	freq.Range.From = time.Date(2017, time.September, 11, 0, 0, 0, 0, time.UTC)
	freq.Target = "100"
	freq.Window.Size = 1

	data, err := json.Marshal(freq)
	suite.NoError(err)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("/api/projects/%s/things/%s/queries/pfetch", projectID, thingID),
		bytes.NewReader(data),
	)
	suite.NoError(err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)

	suite.Equal(200, w.Code)

	suite.NoError(json.Unmarshal(w.Body.Bytes(), &results))

	suite.Equal(1, len(results))
	suite.Equal(6750.0, results[0].Data)
}

func (suite *DMTestSuite) Test_QueriesHandler_Fetch() {
	var results []types.State

	var freq fetchReq
	freq.Range.To = time.Date(2019, time.September, 11, 0, 0, 0, 0, time.UTC)
	freq.Range.From = time.Date(2017, time.September, 11, 0, 0, 0, 0, time.UTC)
	freq.Target = "101"
	freq.Type = "string"

	data, err := json.Marshal(freq)
	suite.NoError(err)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("/api/projects/%s/things/%s/queries/fetch", projectID, thingID),
		bytes.NewReader(data),
	)
	suite.NoError(err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)

	suite.Equal(200, w.Code)

	suite.NoError(json.Unmarshal(w.Body.Bytes(), &results))

	suite.Equal(1, len(results))
	suite.Equal("hello", results[0].Value.String)
}

func (suite *DMTestSuite) Test_QueriesHandler_Recently() {
	var results []types.State

	var rreq recentlyReq
	rreq.Asset = "102"
	rreq.Limit = 1

	data, err := json.Marshal(rreq)
	suite.NoError(err)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("/api/projects/%s/things/%s/queries/recently", projectID, thingID),
		bytes.NewReader(data),
	)
	suite.NoError(err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)

	suite.Equal(200, w.Code)

	suite.NoError(json.Unmarshal(w.Body.Bytes(), &results))

	suite.Equal(1, len(results))
	suite.Equal(7100.0, results[0].Value.Number)
}
