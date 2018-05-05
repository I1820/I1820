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

	"github.com/aiotrc/pm/client"
)

func TestCreate(t *testing.T) {
	h := handle()
	s := httptest.NewServer(h)
	defer s.Close()

	setupDB()

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
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Response: %s", data)
	}
}

func TestThingCreate(t *testing.T) {
	h := handle()
	s := httptest.NewServer(h)
	defer s.Close()

	setupDB()

	d, err := json.Marshal(thingReq{Name: "Me"})
	if err != nil {
		t.Fatalf("Thing request marshaling: %s\n", d)
	}

	resp, err := http.Post(fmt.Sprintf("%s/api/project/%s/things", s.URL, "Her"), "application/json", bytes.NewBuffer(d))
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			return
		}
	}()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Response: %s", data)
	}
}

func TestClient(t *testing.T) {
	h := handle()
	s := httptest.NewServer(h)
	defer s.Close()

	setupDB()

	p := client.New(s.URL)

	pr, err := p.GetThingProject("Me")

	if err != nil {
		t.Fatalf("GetThing error: %s\n", err)
	}

	t.Logf("http://somewhere:%s\n", pr.Runner.Port)
}

func TestThingRemove(t *testing.T) {
	h := handle()
	s := httptest.NewServer(h)
	defer s.Close()

	setupDB()

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/things/%s", s.URL, "Me"), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			return
		}
	}()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Response: %s", data)
	}

}

func TestRemove(t *testing.T) {
	h := handle()
	s := httptest.NewServer(h)
	defer s.Close()

	setupDB()

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
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Response: %s", data)
	}
}
