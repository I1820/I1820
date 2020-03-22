package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/labstack/echo/v4"
)

func (suite *Suite) TestRunnerHandler() {
	suite.testProjectsHandlerCreate()

	// wait for ElRunner make ready
	time.Sleep(15 * time.Second)

	suite.testRunnerAboutAPI()

	suite.testProjectsHandlerDelete()
}

// Runner About API (GET /api/runners/{project_id}/about)
func (suite *Suite) testRunnerAboutAPI() {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", fmt.Sprintf("/runners/%s/about", suite.pID), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
}
