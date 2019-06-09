/*
 * +===============================================
 * | Creation Date: 07-07-2018
 * |
 * | File Name:     thing.go
 * +===============================================
 */

package actions

import (
	"fmt"
	"net/http"

	"github.com/I1820/tm/models"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

// ThingsHandler handles existing things
type ThingsHandler struct {
	db *mongo.Database
}

// thing request payload
type thingReq struct {
	Name     string `json:"name" validate:"required"`
	Model    string `json:"model" validate:"omitempty,alphanum"`
	Location struct {
		Latitude  float64 `json:"lat"`
		Longitude float64 `json:"long"`
	} `json:"location"`
}

// List gets all things. This function is mapped to the path
// GET /projects/{project_id}/things
func (v ThingsHandler) List(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	projectID := c.Param("project_id")

	results := make([]models.Thing, 0)

	cur, err := v.db.Collection("things").Find(ctx, bson.M{
		"project": projectID,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	for cur.Next(ctx) {
		var result models.Thing

		if err := cur.Decode(&result); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		results = append(results, result)
	}
	if err := cur.Close(ctx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, results)
}

// Create adds a thing to the DB and its project. This function is mapped to the
// path POST /projects/{project_id}/things
func (v ThingsHandler) Create(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	projectID := c.Param("project_id")

	var rq thingReq
	if err := c.Bind(&rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// read more about thing model in I1820 platform website
	model := "generic"
	if rq.Model != "" {
		model = rq.Model
	}

	// there is no check for project existence!
	// but sajjad checks it

	t := models.Thing{
		Name:   rq.Name,
		Model:  model,
		Status: true,

		Project: projectID,
	}

	// set thing location if it is provided by user
	// otherwise it would be 0, 0
	/*
		t.Location.Type = "Point"
		t.Location.Coordinates = []float64{rq.Location.Longitude, rq.Location.Latitude}
	*/

	if _, err := v.db.Collection("things").InsertOne(ctx, t); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, t)
}

// Show gets the data for one thing. This function is mapped to
// the path GET /things/{thing_id}
func (v ThingsHandler) Show(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	id := c.Param("thing_id")

	var t models.Thing

	dr := v.db.Collection("things").FindOne(ctx, bson.M{
		"status": true,
		"name":   id,
	})

	if err := dr.Decode(&t); err != nil {
		if err == mongo.ErrNoDocuments {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("thing %s not found", id))
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, t)
}

// Update updates a thing information includes name, model and location. Please note that you must
// provide them all in update request even if you do not want to change it.
// This function is mapped to the path PUT /things/{thing_id}
func (v ThingsHandler) Update(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	id := c.Param("thing_id")

	var rq thingReq
	if err := c.Bind(&rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	dr := v.db.Collection("things").FindOneAndUpdate(ctx, bson.M{
		"name": id,
	}, bson.M{
		"$set": bson.M{
			"model": rq.Model,
		},
	}, options.FindOneAndUpdate().SetReturnDocument(options.After))

	var t models.Thing

	if err := dr.Decode(&t); err != nil {
		if err == mongo.ErrNoDocuments {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("thing %s not found", id))
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, t)

}

// Destroy deletes a thing from the DB and its project. This function is mapped
// to the path DELETE /things/{thing_id}
func (v ThingsHandler) Destroy(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	id := c.Param("thing_id")

	if _, err := v.db.Collection("things").DeleteOne(ctx, bson.M{
		"name": id,
	}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, true)
}

// Activation activates/deactivates thing. This function is mapped
// to the path GET /things/{thing_id}/{t:(?:activate|deactivate)}
func (v ThingsHandler) Activation(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	id := c.Param("thing_id")

	status := false
	if c.Param("t") == "activate" {
		status = true
	}

	dr := v.db.Collection("things").FindOneAndUpdate(ctx, bson.M{
		"name": id,
	}, bson.M{
		"$set": bson.M{"status": status},
	}, options.FindOneAndUpdate().SetReturnDocument(options.After))

	var t models.Thing

	if err := dr.Decode(&t); err != nil {
		if err == mongo.ErrNoDocuments {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("thing %s not found", id))
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, t)
}
