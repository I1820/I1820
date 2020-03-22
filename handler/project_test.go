package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/I1820/I1820/model"
	"github.com/I1820/I1820/request"
	"github.com/labstack/echo/v4"
)

const pName = "raha"
const pOwner = "elahe.dstn@gmail.com"

func (suite *Suite) TestProjectsHandler() {
	suite.testProjectsHandlerCreate()
	suite.testProjectsHandlerShow()
	suite.testProjectsHandlerList(1)
	suite.testProjectsHandlerUpdate()
	suite.testProjectsHandlerDelete()
	suite.testProjectsHandlerList(0)
}

func (suite *Suite) testProjectsHandlerCreate() {
	var p model.Project

	// build project creation request
	var preq request.Project
	preq.Name = pName
	preq.Owner = pOwner
	data, err := json.Marshal(preq)
	suite.NoError(err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/projects", bytes.NewReader(data))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	suite.NoError(json.Unmarshal(w.Body.Bytes(), &p))
	suite.Equal(pName, p.Name)

	suite.pID = p.ID
}

func (suite *Suite) testProjectsHandlerUpdate() {
	var p model.Project

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", fmt.Sprintf("/projects/%s", suite.pID), strings.NewReader("elahe"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	suite.NoError(json.Unmarshal(w.Body.Bytes(), &p))
	suite.Equal("elahe", p.Name)
}

// List GET (/api/projects)
func (suite *Suite) testProjectsHandlerList(count int) {
	var ps []model.Project

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/projects", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	suite.NoError(json.Unmarshal(w.Body.Bytes(), &ps))

	suite.Equal(count, len(ps))
}

// Show GET (/api/projects/{project_id})
func (suite *Suite) testProjectsHandlerShow() {
	var p model.Project

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", fmt.Sprintf("/projects/%s", suite.pID), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	suite.NoError(json.Unmarshal(w.Body.Bytes(), &p))

	suite.Equal(pName, p.Name)
}

// Destroy (DELETE /api/projects/{project_id})
func (suite *Suite) testProjectsHandlerDelete() {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/projects/%s", pName), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
}
