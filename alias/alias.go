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
type Alias struct {
	Project string
	Map     map[string]string
}

// Add adds new mapping between key and name
func (a *Alias) Add(key string, name string) {
	a.Map[key] = name
}
