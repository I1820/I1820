/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 01-02-2019
 * |
 * | File Name:     home_test.go
 * +===============================================
 */

package actions

import (
	"net/http"
	"net/http/httptest"
)

func (suite *WFTestSuite) Test_AboutHandler() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/about", nil)
	suite.NoError(err)
	suite.engine.ServeHTTP(w, req)

	suite.Equal(200, w.Code)
	suite.Contains(w.Body.String(), "18.20 is leaving us")
}
