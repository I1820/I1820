package actions

import (
	"fmt"
	"net/http"

	"github.com/I1820/types"
	"github.com/labstack/echo/v4"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// AssetsHandler handles existing assets
type AssetsHandler struct {
	db *mongo.Database
}

// asset request payload
type assetReq struct {
	Name  string `json:"name" validate:"alphanum,required"`
	Title string `json:"title" validate:"required"`
	Type  string `json:"type" validate:"required,oneof=boolean number string array object location"`
	Kind  string `json:"kind" validate:"required,oneof=sensor actuator"`
}

// List gets all assets of a given thing. This function is mapped to the path
// GET /projects/{project_id}/things/{thing_id}/assets
func (v AssetsHandler) List(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	thingID := c.Param("thing_id")
	projectID := c.Param("project_id")
	var t types.Thing

	dr := v.db.Collection("things").FindOne(ctx, primitive.M{
		"_id":     thingID,
		"project": projectID,
	})
	if err := dr.Decode(&t); err != nil {
		if err == mongo.ErrNoDocuments {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Thing %s not found", thingID))
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, t.Assets)
}
