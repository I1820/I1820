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

package handler

import (
	"testing"

	"github.com/I1820/dm/config"
	"github.com/I1820/dm/db"
	"github.com/I1820/dm/router"
	"github.com/I1820/dm/store"
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
	cfg := config.New()
	db, err := db.New(cfg.Database.URL, "i1820")
	suite.NoError(err)

	suite.engine = router.App(true, "i1820_dm")

	qh := QueriesHandler{
		Store: store.Data{
			DB: db,
		},
	}
	qh.Register(suite.engine.Group(""))
	suite.engine.GET("/about", AboutHandler)
}

// Let's test dm APIs!
func TestService(t *testing.T) {
	suite.Run(t, new(DMTestSuite))
}
