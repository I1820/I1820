/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 07-07-2018
 * |
 * | File Name:     thing_test.go
 * +===============================================
 */

package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/I1820/I1820/request"
	"github.com/labstack/echo/v4"

	"github.com/I1820/I1820/model"
)

const tName = "0000000000000073"
const pID = "raha"

func (suite *TMTestSuite) TestThingsHandler() {
	suite.testThingsHandlerCreate()
	suite.testThingsHandlerShow()
	suite.testThingsHandlerList(1)
	suite.testThingsHandlerDestroy()
	suite.testThingsHandlerShow404()
	suite.testThingsHandlerList(0)
}

// Create thing (POST /api/projects/{project_id}/things)
func (suite *TMTestSuite) testThingsHandlerCreate() {
	var t model.Thing

	// build thing creation request
	var treq request.Thing
	treq.Name = tName
	treq.Location.Latitude = 35.807657 // I1820 location in velenjak
	treq.Location.Longitude = 51.398408
	data, err := json.Marshal(treq)
	suite.NoError(err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", fmt.Sprintf("/projects/%s/things", pID), bytes.NewReader(data))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	suite.NoError(json.Unmarshal(w.Body.Bytes(), &t))
	suite.Equal(tName, t.Name)
}

// Show (GET /api/things/{thing_id}
func (suite *TMTestSuite) testThingsHandlerShow() {
	var t model.Thing

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", fmt.Sprintf("/things/%s", tName), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	suite.NoError(json.Unmarshal(w.Body.Bytes(), &t))

	suite.Equal(tName, t.Name)
}

// List (GET /api/projects/{project_id}/things)
func (suite *TMTestSuite) testThingsHandlerList(count int) {
	var ts []model.Thing

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", fmt.Sprintf("/projects/%s/things", pID), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	suite.NoError(json.Unmarshal(w.Body.Bytes(), &ts))

	suite.Equal(count, len(ts))
}

// Destroy (DELETE /api/things/{thing_id})
func (suite *TMTestSuite) testThingsHandlerDestroy() {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/things/%s", tName), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
}

// Show 404 (GET /api/things/{thing_id}
func (suite *TMTestSuite) testThingsHandlerShow404() {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", fmt.Sprintf("/things/%s", tName), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)

	suite.Equal(http.StatusNotFound, w.Code)
}
