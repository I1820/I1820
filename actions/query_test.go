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

const thingID = "0000000000000073"
const projectID = "el-project"

func (suite *DMTestSuite) Test_QueriesHandler_List() {
	var results []listResp

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/projects/%s/queries/list", projectID), nil)
	suite.NoError(err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)

	suite.Equal(200, w.Code)

	suite.NoError(json.Unmarshal(w.Body.Bytes(), &results))

	suite.NotEqual(0, len(results))

	for _, r := range results {
		if r.ID == thingID {
			suite.Equal(3, r.Total)
		}
	}
}

func (suite *DMTestSuite) Test_QueriesHandler_Fetch() {
	var results []types.Data

	var freq fetchReq
	freq.Since = time.Date(2017, time.September, 11, 0, 0, 0, 0, time.UTC).Unix()
	freq.Until = time.Date(2019, time.September, 11, 0, 0, 0, 0, time.UTC).Unix()
	freq.ThingIDs = []string{thingID}

	data, err := json.Marshal(freq)
	suite.NoError(err)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(
		"POST",
		"/api/queries/fetch",
		bytes.NewReader(data),
	)
	suite.NoError(err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)

	suite.Equal(200, w.Code)

	suite.NoError(json.Unmarshal(w.Body.Bytes(), &results))

	suite.Equal(1, len(results))

	record := results[0]
	suite.Equal(thingID, record.ThingID)
	suite.Equal(7000, record.Data.(map[string]interface{})["100"])
	suite.Equal(6606, record.Data.(map[string]interface{})["101"])
	suite.Equal("17", record.Data.(map[string]interface{})["count"])
}
