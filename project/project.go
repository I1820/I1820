/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 18-11-2017
 * |
 * | File Name:     project.go
 * +===============================================
 */

package project

import "github.com/aiotrc/pm/runner"

// Project represents structure of ISRC projects
type Project struct {
	Name   string        `json:"name"`
	Runner runner.Runner `json:"runner"`
}

// New creates new project with given name
func New(name string, mgu string) (*Project, error) {
	r, err := runner.New(name, mgu)

	if err != nil {
		return nil, err
	}

	return &Project{
		Name:   name,
		Runner: r,
	}, nil
}
