/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 18-11-2017
 * |
 * | File Name:     project.go
 * +===============================================
 */

package model

import (
	"math/rand"
	"time"

	"github.com/I1820/I1820/runner"
)

// Project represents structure of I1820 platform projects
// The project is a virtual entity that collects things together
// under one name and eases their management,
// like an Agricultural project that manages your farm and its smart things.
type Project struct {
	ID     string        `json:"id" bson:"id"`         // Project unique identifier
	Name   string        `json:"name" bson:"name"`     // project human readable name
	Runner runner.Runner `json:"runner" bson:"runner"` // information about project docker

	Description string `json:"description" bson:"description"` // project description

	Perimeter struct { // operational perimeter
		Type        string        `json:"type" bson:"type"`               // GeoJSON type eg. "Polygon"
		Coordinates [][][]float64 `json:"coordinates" bson:"coordinates"` // coordinates eg. [ [ [ 0 , 0 ] , [ 3 , 6 ] , [ 6 , 1 ] , [ 0 , 0  ] ] ]
	} `json:"perimeter" bson:"perimeter"`

	Inspects interface{} `json:"inspects,omitempty" bson:"-"` // more information about project docker
}

// NewProjectID generates a random string as a project identification
func NewProjectID() string {
	rand.Seed(time.Now().UnixNano())

	// Length is a random key length
	const Length = 6

	const source = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

	// Key generates a random key from the source
	b := make([]byte, Length)
	for i := range b {
		b[i] = source[rand.Intn(len(source))]
	}

	return string(b)
}
