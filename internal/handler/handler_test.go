package handler_test

import (
	"testing"

	"github.com/I1820/I1820/internal/config"
	"github.com/I1820/I1820/internal/db"
	"github.com/I1820/I1820/internal/handler"
	"github.com/I1820/I1820/internal/router"
	"github.com/I1820/I1820/internal/store"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

// Suite is a test suite for APIs.
type Suite struct {
	suite.Suite
	engine *echo.Echo

	pID string
}

// SetupSuite initiates tm test suite.
func (suite *Suite) SetupSuite() {
	cfg := config.New()

	db, err := db.New(cfg.Database)
	suite.NoError(err)

	suite.engine = router.App()

	th := handler.Things{
		Store: store.Thing{
			DB: db,
		},
	}
	th.Register(suite.engine.Group(""))

	qh := handler.Queries{
		Store: store.Data{
			DB: db,
		},
	}
	qh.Register(suite.engine.Group(""))

	manager := &handler.MockedManager{}

	ph := handler.Projects{
		Store: store.Project{
			DB: db,
		},
		Manager: manager,
		Config:  cfg.Docker.Runner,
	}
	ph.Register(suite.engine.Group(""))

	rh := handler.Runner{
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
