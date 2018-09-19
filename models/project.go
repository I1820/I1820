/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 18-11-2017
 * |
 * | File Name:     project.go
 * +===============================================
 */

package models

import (
	"context"

	"github.com/I1820/pm/runner"
)

// Project represents structure of I1820 platform projects
// The project is a virtual entity that collects things together
// under one name and eases their management,
// like an Agricultural project that manages your farm and its smart things.
type Project struct {
	ID     string        `json:"id" bson:"_id,omitempty"` // Thing unique identifier
	Name   string        `json:"name" bson:"name"`        // project human readable name
	Runner runner.Runner `json:"runner" bson:"runner"`    // information about project docker

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
