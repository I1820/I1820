package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/I1820/types"
	"github.com/labstack/echo/v4"
)

const aName = "101"
const aTitle = "Temperature"
const aKind = "sensor"
const aType = "number"

func (suite *TMTestSuite) Test_AssetsHandler() {
	suite.testThingsHandlerCreate()

	suite.testAssetsHandlerCreate()
	suite.testAssetsHandlerShow()
	suite.testAssetsHandlerDestroy()

	suite.testThingsHandlerDestroy()
}

// Create asset (POST /api/projects/{project_id}/things/{thing_id}/assets)
func (suite *TMTestSuite) testAssetsHandlerCreate() {
	// build asset creation request
	var areq assetReq
	areq.Name = aName
	areq.Title = aTitle
	areq.Type = aType
	areq.Kind = aKind
	data, err := json.Marshal(areq)
	suite.NoError(err)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", fmt.Sprintf("/api/projects/%s/things/%s/assets", pID, tID), bytes.NewReader(data))
	suite.NoError(err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)

	suite.Equal(200, w.Code)

	var t types.Thing
	suite.NoError(json.Unmarshal(w.Body.Bytes(), &t))
}

// Show (GET /api/projects/{project_id}/things/{thing_id}/assets/{asset_id})
func (suite *TMTestSuite) testAssetsHandlerShow() {
	var a types.Asset

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/projects/%s/things/%s/assets/%s", pID, tID, aName), nil)
	suite.NoError(err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)

	suite.Equal(200, w.Code)

	suite.NoError(json.Unmarshal(w.Body.Bytes(), &a))

	suite.Equal(aTitle, a.Title)
	suite.Equal(aType, a.Type)
	suite.Equal(aKind, a.Kind)
}

// Destroy (DELETE /api/projects/{project_id}/things/assets/{asset_id})
func (suite *TMTestSuite) testAssetsHandlerDestroy() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/api/projects/%s/things/%s/assets/%s", pID, tID, aName), nil)
	suite.NoError(err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)

	suite.Equal(200, w.Code)
}
