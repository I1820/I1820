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

	"github.com/I1820/tm/models"
	"github.com/go-resty/resty"
	"github.com/patrickmn/go-cache"
)

// Error represents project manager errors
type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e Error) Error() string {
	return fmt.Sprintf("(%d): %s", e.Code, e.Message)
}

// TMClient is way for connecting with tm service
type TMService struct {
	cli *resty.Client
	c   *cache.Cache
}

// New creates new instance of PM but connection establishment
// does not happen here.
func New(url string) TMService {
	cli := resty.New().
		SetHostURL(url).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetError(Error{}).
		SetCloseConnection(true)

	return TMService{
		cli: cli,
		c:   cache.New(5*time.Minute, 10*time.Minute),
	}
}

// List lists existing things
func (tm TMService) List() ([]models.Thing, error) {
	var ts []models.Thing

	resp, err := tm.cli.R().
		SetResult(&ts).
		Get("/api/things")
	if err != nil {
		return ts, err
	}

	if resp.StatusCode() != http.StatusOK {
		return ts, resp.Error().(*Error)
	}

	return ts, nil

}

// Show shows thing information by name
func (tm TMService) Show(name string) (models.Thing, error) {
	// check cache first
	if t, found := tm.c.Get(name); found {
		return t.(models.Thing), nil
	}

	var t models.Thing
	resp, err := tm.cli.R().
		SetResult(&t).
		SetPathParams(map[string]string{
			"id": name,
		}).
		Get("/api/things/{id}")
	if err != nil {
		return t, err
	}

	if resp.StatusCode() != http.StatusOK {
		return t, resp.Error().(*Error)
	}

	if !t.Status {
		return t, fmt.Errorf("thing (%s) is not active", name)
	}
	tm.c.Set(t.Name, t, cache.DefaultExpiration)

	return t, nil
}

// Delete deletes thing by name
func (tm TMService) Delete(name string) (models.Thing, error) {
	var t models.Thing

	resp, err := tm.cli.R().
		SetResult(&t).
		SetPathParams(map[string]string{
			"id": name,
		}).
		Delete("/api/things/{id}")
	if err != nil {
		return t, err
	}

	if resp.StatusCode() != http.StatusOK {
		return t, resp.Error().(*Error)
	}

	return t, nil
}
