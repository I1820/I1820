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

	"github.com/I1820/types"
)

const tName = "0000000000000073"
const pID = "kj"

var tID = ""
var t types.Thing

func (suite *TMTestSuite) Test_ThingsHandler() {
	suite.Test_ThingsHandler_Create()
	suite.Test_ThingsHandler_Show()
	suite.Test_ThingsHandler_List()
	suite.Test_ThingsHandler_GeoWithin()
	suite.Test_ThingsHandler_Destroy()
}

// Create thing (POST /api/projects/{project_id}/things)
func (suite *TMTestSuite) Test_ThingsHandler_Create() {
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
	tID = t.ID
}

// Show (GET /api/projects/{project_id}/things/{thing_id}
func (suite *TMTestSuite) Test_ThingsHandler_Show() {
	var ts types.Thing

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/projects/%s/things/%s", pID, tID), nil)
	suite.NoError(err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)

	suite.Equal(200, w.Code)

	suite.NoError(json.Unmarshal(w.Body.Bytes(), &ts))

	suite.Equal(ts, t)
}

// List (GET /api/projects/{project_id}/things)
func (suite *TMTestSuite) Test_ThingsHandler_List() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/projects/%s/things", pID), nil)
	suite.NoError(err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)

	suite.Equal(200, w.Code)
}

// GeoWithin (POST /api/projects/{project_id}/things/geo)
func (suite *TMTestSuite) Test_ThingsHandler_GeoWithin() {
	var tg []types.Thing

	// build thing geowithin request
	var greq = geoWithinReq{
		[][]float64{
			[]float64{35.806731, 51.398618},
			[]float64{35.807784, 51.397810},
			[]float64{35.807827, 51.399516},
			[]float64{35.806731, 51.398618},
		},
	}
	data, err := json.Marshal(greq)
	suite.NoError(err)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", fmt.Sprintf("/api/projects/%s/things/geo", pID), bytes.NewReader(data))
	suite.NoError(err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)

	suite.Equal(200, w.Code)

	suite.NoError(json.Unmarshal(w.Body.Bytes(), &tg))

	suite.Equal(len(tg), 1) // one thing must be found in given location
	suite.Equal(tg[0], t)

}

// Destroy (DELETE /api/projects/{project_id}/things)
func (suite *TMTestSuite) Test_ThingsHandler_Destroy() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/api/projects/%s/things/%s", pID, tID), nil)
	suite.NoError(err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)

	suite.Equal(200, w.Code)
}
