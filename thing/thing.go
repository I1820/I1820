/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 07-02-2018
 * |
 * | File Name:     thing/thing.go
 * +===============================================
 */

package thing

import "github.com/aiotrc/pm/project"

// Thing contains identification and parent project
type Thing struct {
	ID      string
	Project *project.Project
}
