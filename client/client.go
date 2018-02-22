/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 07-02-2018
 * |
 * | File Name:     client/client.go
 * +===============================================
 */

package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/aiotrc/pm/thing"
)

var cache map[string]entry

type entry struct {
	th thing.Thing
	ti time.Time
}

func init() {
	cache = make(map[string]entry)
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
	if e, ok := cache[name]; ok {
		if time.Now().Sub(e.ti) < time.Second {
			return e.th, nil
		}
	}

	var t thing.Thing

	resp, err := http.Get(fmt.Sprintf("%s/api/things/%s", p.URL, name))
	if err != nil {
		return t, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return t, err
	}
	if err := resp.Body.Close(); err != nil {
		return t, err
	}

	if resp.StatusCode != http.StatusOK {
		var e struct {
			Error string
		}

		if err := json.Unmarshal(data, &e); err != nil {
			return t, fmt.Errorf("%s", data)
		}

		return t, fmt.Errorf("%s", e.Error)
	}

	if err := json.Unmarshal(data, &t); err != nil {
		return t, err
	}

	cache[name] = entry{
		th: t,
		ti: time.Now(),
	}

	return t, nil
}
