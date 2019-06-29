package handler

import (
	"os"
	"testing"

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
	mongo := os.Getenv("I1820_TM_DATABASE_URL")
	if mongo == "" {
		mongo = "mongodb://127.0.0.1:27017"
	}
	db, err := db.New(mongo, "i1820")
	suite.NoError(err)

	suite.engine = router.App(true, "i1820_tm")

	th := ThingsHandler{
		Store: store.Things{
			DB: db,
		},
	}
	th.Register(suite.engine.Group("/"))
	suite.engine.GET("/about", AboutHandler)
}

// Let's test tm APIs!
func TestService(t *testing.T) {
	suite.Run(t, new(TMTestSuite))
}
