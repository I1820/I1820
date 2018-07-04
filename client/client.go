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

	"github.com/aiotrc/pm/models"
)

var cache map[string]entry

type entry struct {
	pr models.Project
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

// GetThingProject gets project contains given thing from pm using http request
func (p PM) GetThingProject(name string) (models.Project, error) {
	for _, e := range cache {
		for _, t := range e.pr.Things {
			if t.ID == name {
				if time.Now().Sub(e.ti) < time.Second {
					return e.pr, nil
				}
			}
		}
	}

	var pr models.Project

	resp, err := http.Get(fmt.Sprintf("%s/api/things/%s", p.URL, name))
	if err != nil {
		return pr, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return pr, err
	}
	if err := resp.Body.Close(); err != nil {
		return pr, err
	}

	if resp.StatusCode != http.StatusOK {
		var e struct {
			Error string
		}

		if err := json.Unmarshal(data, &e); err != nil {
			return pr, fmt.Errorf("%s", data)
		}

		return pr, fmt.Errorf("%s", e.Error)
	}

	if err := json.Unmarshal(data, &pr); err != nil {
		return pr, err
	}

	cache[name] = entry{
		pr: pr,
		ti: time.Now(),
	}

	return pr, nil
}
