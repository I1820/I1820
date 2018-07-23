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
	"fmt"
	"net/http"
	"time"

	"github.com/aiotrc/pm/models"
	"github.com/go-resty/resty"
	"github.com/patrickmn/go-cache"
)

var c *cache.Cache

type entry struct {
	pr models.Project
	ti time.Time
}

// Error represents pm errors
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

// ThingsShow shows project that contains given thing
func (p PM) ThingsShow(name string) (models.Project, error) {
	if pr, found := c.Get(name); found {
		return pr.(models.Project), nil
	}

	var pr models.Project

	resp, err := p.cli.R().
		SetResult(&pr).
		SetPathParams(map[string]string{
			"thingId": name,
		}).
		Get("/api/things/{thingId}")
	if err != nil {
		return pr, err
	}

	if resp.StatusCode() != http.StatusOK {
		return pr, resp.Error().(*Error)
	}

	status := false
	for _, t := range pr.Things {
		if t.Status {
			c.Set(t.ID, pr, cache.DefaultExpiration)
			if t.ID == name {
				status = true
			}
		}
	}

	if status {
		return pr, nil
	}

	return pr, fmt.Errorf("Thing (%s) is not activated", name)
}

// RunnerDecode decodes given data on given runner
func (p PM) RunnerDecode(payload []byte, project string, device string) (string, error) {
	resp, err := p.cli.R().
		SetBody(payload).
		SetPathParams(map[string]string{
			"projectId": project,
			"thingId":   device,
		}).
		Post("/api/runners/{projectId}/decode/{thingId}")
	if err != nil {
		return "", err
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("%s", resp.String())
	}

	return resp.String(), nil
}
