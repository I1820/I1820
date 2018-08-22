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

// Project represents structure of ISRC projects
type Project struct {
	Name   string        `json:"name" bson:"name"`
	Runner runner.Runner `json:"runner" bson:"runner"`
	Things []Thing       `json:"things" bson:"things"`
	Status bool          `json:"status" bson:"status"` // active/inactive

	Inspects interface{} `json:"inspects,omitempty" bson:"-"`
}

// NewProject creates new project with given name
func NewProject(ctx context.Context, name string, envs []runner.Env) (*Project, error) {
	r, err := runner.New(ctx, name, envs)

	if err != nil {
		return nil, err
	}

	return &Project{
		Name:   name,
		Runner: r,
		Things: make([]Thing, 0),
		Status: true,
	}, nil
}
