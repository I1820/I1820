package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HealthzHandler shows server is up and running
func HealthzHandler(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}
