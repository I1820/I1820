/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 07-02-2018
 * |
 * | File Name:     main_test.go
 * +===============================================
 */

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreate(t *testing.T) {
	h := handle()
	s := httptest.NewServer(h)
	defer s.Close()

	d, err := json.Marshal(projectReq{Name: "Her"})
	if err != nil {
		t.Fatalf("Project request marshaling: %s\n", d)
	}

	resp, err := http.Post(fmt.Sprintf("%s/api/project", s.URL), "application/json", bytes.NewBuffer(d))
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			return
		}
	}()
	if _, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Fatal(err)
	}
}

func TestRemove(t *testing.T) {
	h := handle()
	s := httptest.NewServer(h)
	defer s.Close()

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/project/%s", s.URL, "Her"), nil)
	cli := http.Client{}

	resp, err := cli.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			return
		}
	}()
	if _, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Fatal(err)
	}
}
