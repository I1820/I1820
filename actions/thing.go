/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 07-07-2018
 * |
 * | File Name:     thing.go
 * +===============================================
 */

package actions

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/I1820/types"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/streadway/amqp"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/segmentio/ksuid"
)

// ThingsHandler handles existing things
type ThingsHandler struct {
	db *mongo.Database

	ch *amqp.Channel
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

// geo within request payload
// each coordinate in coordinates have following standard format
// [latitude, longitude]
type geoWithinReq struct {
	Coordinates [][]float64 `json:"coordinates" validate:"required"`
}

// have tag request payload
type haveTagReq struct {
	Tags []string `json:"tags" validate:"required"`
}

// List gets all things. This function is mapped to the path
// GET /projects/{project_id}/things
func (v ThingsHandler) List(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	projectID := c.Param("project_id")

	results := make([]types.Thing, 0)

	cur, err := v.db.Collection("things").Find(ctx, primitive.M{
		"project": projectID,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	for cur.Next(ctx) {
		var result types.Thing

		if err := cur.Decode(&result); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		results = append(results, result)
	}
	if err := cur.Close(ctx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
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
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := c.Validate(rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// read more about thing model in I1820 platform website
	model := "generic"
	if rq.Model != "" {
		model = rq.Model
	}

	// there is no check for project existence!

	t := types.Thing{
		ID:             primitive.NewObjectID().Hex(),
		Name:           rq.Name,
		Model:          model,
		Status:         true,
		Tokens:         []string{ksuid.New().String()},
		Assets:         make(map[string]types.Asset),
		Connectivities: make(map[string]interface{}),
		Tags:           make([]string, 0),

		Project: projectID,
	}

	// set thing location if it is provided by user
	// otherwise it would be 0, 0
	t.Location.Type = "Point"
	t.Location.Coordinates = []float64{rq.Location.Longitude, rq.Location.Latitude}

	if _, err := v.db.Collection("things").InsertOne(ctx, t); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	go func() {
		b, err := json.Marshal(t)
		if err != nil {
			log.Errorf("Thing marshaling error: %s", err)
		}

		// publish thing cration event into the direct exchange of rabbitmq
		if err := v.ch.Publish(
			"i1820_things", // exchange type
			"create",       // routing key
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        b,
			},
		); err != nil {
			log.Errorf("Rabbitmq failed to publish: %s", err)
		}
	}()

	return c.JSON(http.StatusOK, t)
}

// Show gets the data for one thing. This function is mapped to
// the path GET /projects/{project_id}/things/{thing_id}
func (v ThingsHandler) Show(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	projectID := c.Param("project_id")
	id := c.Param("thing_id")

	var t types.Thing

	dr := v.db.Collection("things").FindOne(ctx, primitive.M{
		"status":  true,
		"_id":     id,
		"project": projectID,
	})

	if err := dr.Decode(&t); err != nil {
		if err == mongo.ErrNoDocuments {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Thing %s not found", id))
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, t)
}

// Update updates a thing information includes name, model and location. Please note that you must
// provide them all in update request even if you do not want to change it.
// This function is mapped to the path PUT /projects/{project_id}/things/{thing_id}
func (v ThingsHandler) Update(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	projectID := c.Param("project_id")
	id := c.Param("thing_id")

	var rq thingReq
	if err := c.Bind(&rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := c.Validate(rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	dr := v.db.Collection("things").FindOneAndUpdate(ctx, primitive.M{
		"_id":     id,
		"project": projectID,
	}, primitive.M{
		"$set": primitive.M{
			"name":  rq.Name,
			"model": rq.Model,
			"location.coordinates": primitive.A{
				rq.Location.Longitude,
				rq.Location.Latitude,
			},
		},
	}, options.FindOneAndUpdate().SetReturnDocument(options.After))

	var t types.Thing

	if err := dr.Decode(&t); err != nil {
		if err == mongo.ErrNoDocuments {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Thing %s not found", id))
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, t)

}

// Destroy deletes a thing from the DB and its project. This function is mapped
// to the path DELETE /projects/{project_id}/things/{thing_id}
func (v ThingsHandler) Destroy(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	projectID := c.Param("project_id")
	id := c.Param("thing_id")

	if _, err := v.db.Collection("things").DeleteOne(ctx, primitive.M{
		"_id":     id,
		"project": projectID,
	}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, true)
}

// Activation activates/deactivates thing. This function is mapped
// to the path GET /projects/{project_id}/things/{thing_id}/{t:(?:activate|deactivate)}
func (v ThingsHandler) Activation(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	id := c.Param("thing_id")
	projectID := c.Param("project_id")

	status := false
	if c.Param("t") == "activate" {
		status = true
	}

	dr := v.db.Collection("things").FindOneAndUpdate(ctx, primitive.M{
		"_id":     id,
		"project": projectID,
	}, primitive.M{
		"$set": primitive.M{"status": status},
	}, options.FindOneAndUpdate().SetReturnDocument(options.After))

	var t types.Thing

	if err := dr.Decode(&t); err != nil {
		if err == mongo.ErrNoDocuments {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("Thing %s not found", id))
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, t)
}

// GeoWithin returns all things that are in the polygon that is given by user.
// This function is mapped to the path POST /projects/{project_id}/things/geo
func (v ThingsHandler) GeoWithin(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	projectID := c.Param("project_id")

	var rq geoWithinReq
	if err := c.Bind(&rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := c.Validate(rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	coordinates := primitive.A{}
	for _, coordinate := range rq.Coordinates {
		coordinates = append(coordinates, primitive.A{
			coordinate[1], // longitude is first in mongo
			coordinate[0], // latitude is second in mongo
		})
	}

	results := make([]types.Thing, 0)

	cur, err := v.db.Collection("things").Find(ctx, primitive.M{
		"project": projectID,
		"location": primitive.M{
			"$geoWithin": primitive.M{
				"$geometry": primitive.M{
					"type":        "Polygon",
					"coordinates": primitive.A{coordinates},
				},
			},
		},
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	for cur.Next(ctx) {
		var result types.Thing

		if err := cur.Decode(&result); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		results = append(results, result)
	}
	if err := cur.Close(ctx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, results)
}

// HaveTags returns all things that have tags that are given by user
// This function is mapped to the path POST /projeects/{project_id}/things/tags
func (v ThingsHandler) HaveTags(c echo.Context) error {
	// gets the request context
	ctx := c.Request().Context()

	projectID := c.Param("project_id")

	var rq haveTagReq
	if err := c.Bind(&rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := c.Validate(rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	results := make([]types.Thing, 0)

	cur, err := v.db.Collection("things").Find(ctx, primitive.M{
		"project": projectID,
		"tags": primitive.M{
			"$in": rq.Tags,
		},
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	for cur.Next(ctx) {
		var result types.Thing

		if err := cur.Decode(&result); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		results = append(results, result)
	}
	if err := cur.Close(ctx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, results)
}
