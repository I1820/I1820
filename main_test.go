/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 28-04-2018
 * |
 * | File Name:     main_test.go
 * +===============================================
 */

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAbout(t *testing.T) {
	h := handle()
	s := httptest.NewServer(h)
	defer s.Close()

	resp, err := http.Get(fmt.Sprintf("%s/api/about", s.URL))
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if err := resp.Body.Close(); err != nil {
		t.Fatal(err)
	}

	if string(body) != "18.20 is leaving us" {
		t.Fatalf("who leaving us?! %q", body)
	}
}

func TestThings(t *testing.T) {
	setupDB()

	h := handle()
	s := httptest.NewServer(h)
	defer s.Close()

	resp, err := http.Get(fmt.Sprintf("%s/api/things", s.URL))
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var results []struct {
		ID    string `json:"_id"`
		Total int
	}
	if err := json.Unmarshal(body, &results); err != nil {
		t.Fatal(err)
	}
	t.Log(results)

	if len(results) == 0 {
		t.Fatal(fmt.Errorf("Invalid number of records: %d", len(results)))
	}

	for _, r := range results {
		if r.ID == "0000000000000003" {
			if r.Total != 1 {
				t.Fatal(fmt.Errorf("Invalid number of thing 0000000000000003 records: %d", r.Total))
			}
		}
	}

	if err := resp.Body.Close(); err != nil {
		t.Fatal(err)
	}

}
