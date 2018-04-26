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
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/brocaar/lora-app-server/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// LoRaServer represents loraserver.io endpoint
type LoRaServer struct {
	BaseURL string

	jwtToken string
	hc       *http.Client

	username string
	password string
}

// New creates new loraserver.io endpoint and connects to it
func New(baseURL, username, password string) (*LoRaServer, error) {
	l := &LoRaServer{
		BaseURL: baseURL,

		hc: &http.Client{
			Transport: &http.Transport{
				// Disable https certificate validation
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},

		username: username,
		password: password,
	}

	return l, l.login()
}

func (l *LoRaServer) login() error {
	d, _ := json.Marshal(map[string]string{
		"username": l.username,
		"password": l.password,
	})
	resp, err := l.hc.Post("https://"+l.BaseURL+"/api/internal/login", "application/json", bytes.NewBuffer(d))
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("StatusCode: %d", resp.StatusCode)
	}

	var response struct {
		Jwt string
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("JSON Unmarshal: %s", err)
	}

	if err := resp.Body.Close(); err != nil {
		return err
	}

	l.jwtToken = response.Jwt
	return nil
}

// GatewayFrameStream streams gateway frame logs
func (l *LoRaServer) GatewayFrameStream(mac string) (<-chan *GatewayFrame, error) {
	grpcOpts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(jwt{
			token: l.jwtToken,
		}),
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})),
	}

	asConn, err := grpc.Dial(l.BaseURL, grpcOpts...)
	if err != nil {
		return nil, err
	}

	gc := api.NewGatewayClient(asConn)
	s, err := gc.StreamFrameLogs(context.Background(), &api.StreamGatewayFrameLogsRequest{
		Mac: mac,
	})
	if err != nil {
		return nil, err
	}

	c := make(chan *GatewayFrame, 1024)

	go func() {
		for {
			d, err := s.Recv()
			if err != nil {
				return
			}

			for _, up := range d.UplinkFrames {
				c <- &GatewayFrame{
					Mac:           mac,
					UplinkFrame:   up,
					DownlinkFrame: nil,

					Timestamp: time.Now(),
				}
			}
			for _, dl := range d.DownlinkFrames {
				c <- &GatewayFrame{
					Mac:           mac,
					UplinkFrame:   nil,
					DownlinkFrame: dl,

					Timestamp: time.Now(),
				}
			}

		}
	}()

	return c, nil
}
