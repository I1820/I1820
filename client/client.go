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
	"net/http"
	"time"

	"github.com/I1820/pm/models"
	"github.com/go-resty/resty"
	"github.com/patrickmn/go-cache"
)

var c *cache.Cache

// Thing is required information for client from thing
// This information is enough for client to communicate
// with pm
type Thing struct {
	Project string
	Model   string
}

// Error represents project manager errors
type Error struct {
	Err  string `json:"error"`
	Code int    `json:"code"`
}

func (e Error) Error() string {
	return fmt.Sprintf("(%d): %s", e.Code, e.Err)
}

func init() {
	c = cache.New(5*time.Minute, 10*time.Minute)
}

// PM is way for connecting to PM :joy:
type PM struct {
	cli *resty.Client
}

// New creates new instance of PM but connection establishment
// does not happen here.
func New(url string) PM {
	cli := resty.New().
		SetHostURL(url).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetError(Error{}).
		SetCloseConnection(true)

	return PM{
		cli: cli,
	}
}

// ProjectsCreate creates new project
func (p PM) ProjectsCreate(name string) (models.Project, error) {
	var pr models.Project

	resp, err := p.cli.R().
		SetBody(map[string]string{"name": name}).
		SetResult(&pr).
		Post("/api/projects")
	if err != nil {
		return pr, err
	}

	if resp.StatusCode() != http.StatusOK {
		return pr, resp.Error().(*Error)
	}

	return pr, nil
}

// ProjectsList lists existing projects
func (p PM) ProjectsList() ([]models.Project, error) {
	var pr []models.Project

	resp, err := p.cli.R().
		SetResult(&pr).
		Get("/api/projects")
	if err != nil {
		return pr, err
	}

	if resp.StatusCode() != http.StatusOK {
		return pr, resp.Error().(*Error)
	}

	return pr, nil

}

// ProjectsShow shows project information by name
func (p PM) ProjectsShow(name string) (models.Project, error) {
	var pr models.Project

	resp, err := p.cli.R().
		SetResult(&pr).
		SetPathParams(map[string]string{
			"projectId": name,
		}).
		Get("/api/projects/{projectId}")
	if err != nil {
		return pr, err
	}

	if resp.StatusCode() != http.StatusOK {
		return pr, resp.Error().(*Error)
	}

	return pr, nil

}

// ProjectsDelete deletes project by name
func (p PM) ProjectsDelete(name string) (models.Project, error) {
	var pr models.Project

	resp, err := p.cli.R().
		SetResult(&pr).
		SetPathParams(map[string]string{
			"projectId": name,
		}).
		Delete("/api/projects/{projectId}")
	if err != nil {
		return pr, err
	}

	if resp.StatusCode() != http.StatusOK {
		return pr, resp.Error().(*Error)
	}

	return pr, nil

}

// ThingsShow shows information of thing and project that contains that
func (p PM) ThingsShow(name string) (Thing, error) {
	if th, found := c.Get(name); found {
		return th.(Thing), nil
	}

	var pr models.Project
	var th Thing

	resp, err := p.cli.R().
		SetResult(&pr).
		SetPathParams(map[string]string{
			"thingId": name,
		}).
		Get("/api/things/{thingId}")
	if err != nil {
		return th, err
	}

	if resp.StatusCode() != http.StatusOK {
		return th, resp.Error().(*Error)
	}

	status := false
	for _, t := range pr.Things {
		if t.Status {
			// set cache for all things of the project
			c.Set(t.ID, Thing{pr.Name, t.Model}, cache.DefaultExpiration)
			if t.ID == name {
				status = true
				th = Thing{pr.Name, t.Model}
			}
		}
	}

	if status {
		return th, nil
	}

	return th, fmt.Errorf("Thing (%s) is not activate", name)
}

// RunnersDecode decodes given data on given runner
func (p PM) RunnersDecode(payload []byte, project string, device string) (interface{}, error) {
	var result interface{}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := p.cli.R().
		SetBody(body).
		SetResult(&result).
		SetPathParams(map[string]string{
			"projectId": project,
			"thingId":   device,
		}).
		Post("/api/runners/{projectId}/api/codecs/{thingId}/decode")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("%s", resp.String())
	}

	return result, nil
}

//
