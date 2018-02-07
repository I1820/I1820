/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 07-02-2018
 * |
 * | File Name:     pm/pm.go
 * +===============================================
 */

package pm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aiotrc/pm/thing"
)

var cache map[string]thing.Thing

func init() {
	cache = make(map[string]thing.Thing)
}

// PM is way for connecting to PM :joy:
type PM struct {
	URL string
}

// New creates new instance of PM but connection establishment
// does not happen here.
func New(url string) PM {
	return PM{
		URL: url,
	}
}

// GetThing gets thing information from pm using http request
func (p PM) GetThing(name string) (thing.Thing, error) {
	if t, ok := cache[name]; ok {
		return t, nil
	}

	var t thing.Thing

	resp, err := http.Get(fmt.Sprintf("%s/api/thing/%s", p.URL, name))
	if err != nil {
		return t, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return t, err
	}

	if resp.StatusCode != http.StatusOK {
		return t, fmt.Errorf("%s", data)
	}

	if err := json.Unmarshal(data, &t); err != nil {
		return t, err
	}

	cache[name] = t

	return t, nil
}
