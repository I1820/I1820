/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 23-04-2018
 * |
 * | File Name:     main.go
 * +===============================================
 */

package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/brocaar/lora-app-server/api"
	"github.com/go-resty/resty"
	"github.com/jinzhu/configor"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// JWT Token
var jwtToken string

// Config represents main configuration
var Config = struct {
	Loraserver struct {
		BaseURL string `default:"https://platform.ceit.aut.ac.ir:50013/api" env:"backback_base_url"`
	}
}{}

type jwt struct {
	token string
}

func (j jwt) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": j.token,
	}, nil
}

func (j jwt) RequireTransportSecurity() bool {
	return true
}

func main() {
	// Disable https certificate validation
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// Load configuration
	if err := configor.Load(&Config, "config.yml"); err != nil {
		panic(err)
	}

	login()

	grpcOpts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(jwt{
			token: jwtToken,
		}),
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})),
	}

	asConn, err := grpc.Dial("platform.ceit.aut.ac.ir:50013", grpcOpts...)
	if err != nil {
		log.WithFields(log.Fields{
			"Phase": "grpc login",
		}).Fatal(err)
	}

	gc := api.NewGatewayClient(asConn)
	s, err := gc.StreamFrameLogs(context.Background(), &api.StreamGatewayFrameLogsRequest{
		Mac: "b827ebffff47d1a5",
	})
	if err != nil {
		log.WithFields(log.Fields{
			"Phase": "grpc stream frame logs",
		}).Fatal(err)
	}

	d, err := s.Recv()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(d)
}

func login() {
	resp, err := resty.R().
		SetBody(map[string]string{
			"username": "admin",
			"password": "admin",
		}).
		Post(Config.Loraserver.BaseURL + "/internal/login")
	if err != nil {
		log.WithFields(log.Fields{
			"Phase": "login",
		}).Fatalf("Request: %s", err)
	}

	if resp.StatusCode() != 200 {
		log.WithFields(log.Fields{
			"Phase": "login",
		}).Fatalf("StatusCode: %d", resp.StatusCode())
	}

	var response struct {
		Jwt string
	}
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		log.WithFields(log.Fields{
			"Phase": "login",
		}).Fatalf("JSON Unmarshal: %s", err)
	}

	log.WithFields(log.Fields{
		"Phase": "login",
	}).Infoln(response)

	jwtToken = response.Jwt
}
