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
	Name string
	ID   string
}

// New creates new project with given name
func New(name string) *Project {
	return &Project{
		Name: name,
		ID:   runner.New(name),
	}
}
