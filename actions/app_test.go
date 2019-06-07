/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 01-02-2019
 * |
 * | File Name:     app_test.go
 * +===============================================
 */

package actions

import (
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

// TMTestSuite is a test suite for tm component APIs.
type TMTestSuite struct {
	suite.Suite
	engine *echo.Echo
}

// SetupSuite initiates tm test suite
func (suite *TMTestSuite) SetupSuite() {
	mongo := os.Getenv("I1820_TM_DATABASE_URL")
	if mongo == "" {
		mongo = "mongodb://127.0.0.1:27017"
	}
	suite.engine = App(true, mongo)
}

// Let's test tm APIs!
func TestService(t *testing.T) {
	suite.Run(t, new(TMTestSuite))
}
