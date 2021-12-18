package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/I1820/I1820/internal/model"
	"github.com/I1820/I1820/internal/request"
	"github.com/labstack/echo/v4"
)

const (
	thingID   = "0000000000000073"
	projectID = "el-project"
)

func (suite *Suite) TestQueriesHandlerList() {
	require := suite.Require()

	var results map[string]int

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", fmt.Sprintf("/queries/projects/%s/list", projectID), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)

	require.Equal(http.StatusOK, w.Code, w.Body.String())
	require.NoError(json.Unmarshal(w.Body.Bytes(), &results))
	require.Equal(results[thingID], 4)
}

func (suite *Suite) TestQueriesHandlerFetch() {
	require := suite.Require()

	var results []model.Data

	var freq request.Fetch
	freq.Since = time.Date(2017, time.September, 11, 0, 0, 0, 0, time.UTC).Unix()
	freq.Until = time.Date(2019, time.September, 11, 0, 0, 0, 0, time.UTC).Unix()
	freq.ThingIDs = []string{thingID}

	data, err := json.Marshal(freq)
	suite.NoError(err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(
		"POST",
		"/queries/fetch",
		bytes.NewReader(data),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)

	require.Equal(http.StatusOK, w.Code, w.Body.String())
	require.NoError(json.Unmarshal(w.Body.Bytes(), &results))
	require.Equal(4, len(results))

	require.Equal(nil, results[0].Data)

	record := results[1]
	require.Equal(thingID, record.ThingID)
	require.Equal(7000.0, record.Data.(map[string]interface{})["100"])
	require.Equal(6606.0, record.Data.(map[string]interface{})["101"])
	require.Equal("19", record.Data.(map[string]interface{})["count"])
}

func (suite *Suite) TestQueriesHandlerFetchSingle() {
	require := suite.Require()

	var results []model.Data

	since := time.Date(2017, time.September, 11, 0, 0, 0, 0, time.UTC).Unix()
	until := time.Date(2019, time.September, 11, 0, 0, 0, 0, time.UTC).Unix()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(
		"GET",
		fmt.Sprintf("/queries/things/%s/fetch?since=%d&until=%d", thingID, since, until),
		nil,
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)

	require.Equal(http.StatusOK, w.Code, w.Body.String())
	require.NoError(json.Unmarshal(w.Body.Bytes(), &results))
	require.Equal(4, len(results))

	require.Equal(nil, results[0].Data)

	record := results[1]
	require.Equal(thingID, record.ThingID)
	require.Equal(7000.0, record.Data.(map[string]interface{})["100"])
	require.Equal(6606.0, record.Data.(map[string]interface{})["101"])
	require.Equal("19", record.Data.(map[string]interface{})["count"])
}
