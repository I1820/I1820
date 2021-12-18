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

// Suite is a test suite for APIs.
type Suite struct {
	suite.Suite
	engine *echo.Echo

	pID string
}

// SetupSuite initiates tm test suite
func (suite *Suite) SetupSuite() {
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

	qh := Queries{
		Store: store.Data{
			DB: db,
		},
	}
	qh.Register(suite.engine.Group(""))

	manager := &MockedManager{}

	ph := Projects{
		Store: store.Project{
			DB: db,
		},
		Manager: manager,
		Config:  cfg.Docker.Runner,
	}
	ph.Register(suite.engine.Group(""))

	rh := Runner{
		Store: store.Project{
			DB: db,
		},
		Manager:    manager,
		DockerHost: cfg.Docker.Host,
	}
	rh.Register(suite.engine.Group(""))
}

// Let's test APIs!
func TestService(t *testing.T) {
	suite.Run(t, new(Suite))
}
