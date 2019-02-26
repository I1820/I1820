package actions

import (
	"fmt"
	"net/http"

	"github.com/I1820/types"
	"github.com/labstack/echo/v4"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
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

// Create adds an asset to the DB and its thing. This function is mapped to the
// path POST /projects/{project_id}/things/{thing_id}/assets
func (v AssetsHandler) Create(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	thingID := c.Param("thing_id")
	projectID := c.Param("project_id")

	var rq assetReq
	if err := c.Bind(&rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	a := types.Asset{
		Title: rq.Title,
		Type:  rq.Type,
		Kind:  rq.Kind,
	}

	dr := v.db.Collection("things").FindOneAndUpdate(ctx, primitive.M{
		"_id":     thingID,
		"project": projectID,
	}, primitive.M{
		"$set": primitive.M{
			fmt.Sprintf("assets.%s", rq.Name): a,
		},
	}, options.FindOneAndUpdate().SetReturnDocument(options.After))

	var t types.Thing

	if err := dr.Decode(&t); err != nil {
		if err == mongo.ErrNoDocuments {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Thing %s not found", thingID))
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, t)
}

// Show gets the data for a given asset. This function is mapped to
// the path GET /projects/{project_id}/things/{thing_id}/assets/{asset_name}
func (v AssetsHandler) Show(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	thingID := c.Param("thing_id")
	projectID := c.Param("project_id")
	assetName := c.Param("asset_id")

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

	return c.JSON(http.StatusOK, t.Assets[assetName])
}

// Destroy deletes an asset from the DB and its thing. This function is mapped
// to the path DELETE /projects/{project_id}/things/{thing_id}/assets/{asset_name}
func (v AssetsHandler) Destroy(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	thingID := c.Param("thing_id")
	projectID := c.Param("project_id")
	assetName := c.Param("asset_id")

	dr := v.db.Collection("things").FindOneAndUpdate(ctx, primitive.M{
		"_id":     thingID,
		"project": projectID,
	}, primitive.M{
		"$unset": primitive.M{
			fmt.Sprintf("assets.%s", assetName): "",
		},
	}, options.FindOneAndUpdate().SetReturnDocument(options.After))

	var t types.Thing

	if err := dr.Decode(&t); err != nil {
		if err == mongo.ErrNoDocuments {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Thing %s not found", thingID))
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, t)
}
