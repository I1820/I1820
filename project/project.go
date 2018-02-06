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
// each project has name and contains one or more things
type Project struct {
	Name   string
	Runner runner.Runner
	Things []string
}

// New creates new project with given name
func New(name string) *Project {
	return &Project{
		Name:   name,
		Runner: runner.New(name),
		Things: make([]string, 0),
	}
}

// AddThing adds new things into specific project
func (p *Project) AddThing(thing string) {
	p.Things = append(p.Things, thing)
}
