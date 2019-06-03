package actions

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/labstack/echo/v4"
)

// LinkTestSuite is a test suite for link component APIs.
type LinkTestSuite struct {
	suite.Suite
	engine *echo.Echo
}

// SetupSuite initiates link test suite
func (suite *LinkTestSuite) SetupSuite() {
	suite.engine = App(true)
}

// Let's test link APIs!
func TestService(t *testing.T) {
	suite.Run(t, new(LinkTestSuite))
}
