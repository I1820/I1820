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
	"context"

	"github.com/I1820/I1820/runner"
)

// Project represents structure of I1820 platform projects
// The project is a virtual entity that collects things together
// under one name and eases their management,
// like an Agricultural project that manages your farm and its smart things.
type Project struct {
	ID     string        `json:"id" bson:"_id,omitempty"` // Project unique identifier
	Name   string        `json:"name" bson:"name"`        // project human readable name
	Runner runner.Runner `json:"runner" bson:"runner"`    // information about project docker

	Description string   `json:"description" bson:"description"` // project description
	Perimeter   struct { // operational perimeter
		Type        string        `json:"type" bson:"type"`               // GeoJSON type eg. "Polygon"
		Coordinates [][][]float64 `json:"coordinates" bson:"coordinates"` // coordinates eg. [ [ [ 0 , 0 ] , [ 3 , 6 ] , [ 6 , 1 ] , [ 0 , 0  ] ] ]
	} `json:"perimeter" bson:"perimeter"`

	Inspects interface{} `json:"inspects,omitempty" bson:"-"` // more information about project docker
}

// NewProject creates new project with given identification
func NewProject(ctx context.Context, id string, name string, envs []runner.Env) (*Project, error) {
	r, err := runner.New(ctx, id, envs)
	if err != nil {
		return nil, err
	}

	return &Project{
		Name:   name,
		Runner: r,
	}, nil
}

// ReProject recreates project dockers and replaces runner field of given project with new
// information
func ReProject(ctx context.Context, envs []runner.Env, p *Project) error {
	r, err := runner.New(ctx, p.ID, envs)
	if err != nil {
		return err
	}

	p.Runner = r

	return nil
}
