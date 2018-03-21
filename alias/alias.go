/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 21-03-2018
 * |
 * | File Name:     alias/alias.go
 * +===============================================
 */

package alias

// Alias provides a structure for storing data fields aliasing
// Name is a thing identification and aliases are defined per things
type Alias struct {
	Name string
	Map  map[string]string
}

// New creates new empty Alias
func New(name string) *Alias {
	return &Alias{
		Name: name,
		Map:  make(map[string]string),
	}
}

// Add adds new mapping between key and name
func (a *Alias) Add(key string, name string) {
	a.Map[key] = name
}
