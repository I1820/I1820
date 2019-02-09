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
	suite.engine = App()
}

// Let's test dm APIs!
func TestService(t *testing.T) {
	suite.Run(t, new(DMTestSuite))
}
