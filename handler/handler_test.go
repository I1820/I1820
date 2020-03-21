package handler

import (
	"testing"

	"github.com/I1820/tm/config"
	"github.com/I1820/tm/db"
	"github.com/I1820/tm/router"
	"github.com/I1820/tm/store"
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
	db, err := db.New(cfg.Database.URL, "i1820")
	suite.NoError(err)

	suite.engine = router.App(true, "i1820_tm")

	th := Things{
		Store: store.Things{
			DB: db,
		},
	}
	th.Register(suite.engine.Group(""))
	suite.engine.GET("/about", AboutHandler)
}

// Let's test tm APIs!
func TestService(t *testing.T) {
	suite.Run(t, new(TMTestSuite))
}
