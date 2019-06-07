/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 07-07-2018
 * |
 * | File Name:     thing_test.go
 * +===============================================
 */

package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"

	"github.com/I1820/tm/models"
)

const tName = "0000000000000073"
const pID = "kj"

var t models.Thing

func (suite *TMTestSuite) Test_ThingsHandler() {
	suite.testThingsHandlerCreate()
	suite.testThingsHandlerShow()
	suite.testThingsHandlerList()
	suite.testThingsHandlerDestroy()
	suite.testThingsHandlerShow404()
	suite.testThingsHandlerList()
}

// Create thing (POST /api/projects/{project_id}/things)
func (suite *TMTestSuite) testThingsHandlerCreate() {
	// build thing creation request
	var treq thingReq
	treq.Name = tName
	treq.Location.Latitude = 35.807657 // I1820 location in velenjak
	treq.Location.Longitude = 51.398408
	data, err := json.Marshal(treq)
	suite.NoError(err)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", fmt.Sprintf("/api/projects/%s/things", pID), bytes.NewReader(data))
	suite.NoError(err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)

	suite.Equal(200, w.Code)

	suite.NoError(json.Unmarshal(w.Body.Bytes(), &t))
	suite.Equal(tName, t.Name)
}

// Show (GET /api/things/{thing_id}
func (suite *TMTestSuite) testThingsHandlerShow() {
	var ts models.Thing

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/things/%s", tName), nil)
	suite.NoError(err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)

	suite.Equal(200, w.Code)

	suite.NoError(json.Unmarshal(w.Body.Bytes(), &ts))

	suite.Equal(ts, t)
}

// List (GET /api/projects/{project_id}/things)
func (suite *TMTestSuite) testThingsHandlerList() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/projects/%s/things", pID), nil)
	suite.NoError(err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)

	suite.Equal(200, w.Code)
}

// Destroy (DELETE /api/things/{thing_id})
func (suite *TMTestSuite) testThingsHandlerDestroy() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/api/things/%s", tName), nil)
	suite.NoError(err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)

	suite.Equal(200, w.Code)
}

// Show 404 (GET /api/things/{thing_id}
func (suite *TMTestSuite) testThingsHandlerShow404() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/things/%s", tName), nil)
	suite.NoError(err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)

	suite.Equal(404, w.Code)
}
