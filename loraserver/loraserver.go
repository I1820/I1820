/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 24-04-2018
 * |
 * | File Name:     loraserver/loraserver.go
 * +===============================================
 */

package loraserver

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// LoRaServer represents loraserver.io endpoint
type LoRaServer struct {
	BaseURL string

	jwtToken string
	hc       *http.Client
}

// New creates new loraserver.io endpoint and connects to it
func New(baseURL string) (*LoRaServer, error) {
	l := &LoRaServer{
		BaseURL: baseURL,

		hc: &http.Client{
			Transport: &http.Transport{
				// Disable https certificate validation
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}

	return l, l.login()
}

func (l *LoRaServer) login() error {
	d, _ := json.Marshal(map[string]string{
		"username": "admin",
		"password": "admin",
	})
	resp, err := l.hc.Post(l.BaseURL+"/internal/login", "application/json", bytes.NewBuffer(d))
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("StatusCode: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	var response struct {
		Jwt string
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("JSON Unmarshal: %s", err)
	}

	l.jwtToken = response.Jwt
	return nil
}
