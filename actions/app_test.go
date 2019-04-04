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

// DMTestSuite is a test suite for dm component APIs.
type DMTestSuite struct {
	suite.Suite
	engine *echo.Echo
}

// SetupSuite initiates dm test suite
func (suite *DMTestSuite) SetupSuite() {
	mongo := os.Getenv("I1820_DM_DATABASE_URL")
	if mongo == "" {
		mongo = "mongodb://127.0.0.1:27017"
	}
	suite.engine = App(mongo, true)
}

// Let's test dm APIs!
func TestService(t *testing.T) {
	suite.Run(t, new(DMTestSuite))
}
