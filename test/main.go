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
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/brocaar/lora-app-server/api"
	"github.com/go-resty/resty"
	"github.com/jinzhu/configor"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"golang.org/x/net/http2"
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
	http.DefaultTransport = &http2.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

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

func gatewayFrames() {
	resp, err := resty.R().
		SetAuthToken(jwtToken).
		SetDoNotParseResponse(false).
		Get(Config.Loraserver.BaseURL + "/gateways/b827ebffff47d1a5/frames")
	if err != nil {
		log.WithFields(log.Fields{
			"Phase": "gatewayFrames",
		}).Fatalf("Request: %s", err)
	}

	defer resp.RawBody().Close()

	framer := http2.NewFramer(ioutil.Discard, resp.RawBody())
	for {
		f, err := framer.ReadFrame()
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
		switch err.(type) {
		case nil:
			log.Println(f)
		case http2.ConnectionError:
			// Ignore. There will be many errors of type "PROTOCOL_ERROR, DATA
			// frame with stream ID 0". Presumably we are abusing the framer.
		default:
			log.Println(err, framer.ErrorDetail())
		}
	}

	if resp.StatusCode() != 200 {
		log.WithFields(log.Fields{
			"Phase": "gatewayFrames",
		}).Fatalf("StatusCode: %d", resp.StatusCode())
	}

	if _, err := io.Copy(os.Stdout, resp.RawBody()); err != nil {
		log.WithFields(log.Fields{
			"Phase": "gatewayFrames",
		}).Fatalf("io.Copy: %s", err)
	}

	var response struct {
		UplinkFrames []struct {
			PhyPayloadJSON string
			Frequency      int
		}
	}

	log.WithFields(log.Fields{
		"Phase": "gatewayFrames",
	}).Infoln(response)
}
