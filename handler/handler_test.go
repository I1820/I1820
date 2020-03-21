package handler

import (
	"testing"

	"github.com/I1820/I1820/config"
	"github.com/I1820/I1820/db"
	"github.com/I1820/I1820/router"
	"github.com/I1820/I1820/store"
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
	cfg := config.New()
	db, err := db.New(cfg.Database)
	suite.NoError(err)

	suite.engine = router.App()

	th := Things{
		Store: store.Thing{
			DB: db,
		},
	}

	th.Register(suite.engine.Group(""))
}

// Let's test tm APIs!
func TestService(t *testing.T) {
	suite.Run(t, new(TMTestSuite))
}
