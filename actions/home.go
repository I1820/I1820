package actions

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// AboutHandler is a default handler to serve up
// memory of 18.20.
func AboutHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "18.20 is leaving us")
}
