package mqtt

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/I1820/I1820/internal/config"
	"github.com/I1820/I1820/internal/model"
	"github.com/I1820/I1820/pkg/protocol"
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

// Service is MQTT service for handling device communication.
type Service struct {
	cli  paho.Client
	opts *paho.ClientOptions

	protocols []protocol.Protocol

	channel chan model.Data

	IsRun bool
}

// MaximumClientID represents maximum client identification of MQTT.
const MaximumClientID = 1024

// DisconnectTimeout is a waiting time for MQTT client disconnect in ms.
const DisconnectTimeout = 10

// New creates a new MQTT service instance.
// MQTT service receives messages from protocols and publish them on its channel.
func New(cfg config.MQTT) *Service {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	opts := paho.NewClientOptions()
	opts.AddBroker(cfg.Addr)
	opts.SetClientID(fmt.Sprintf("I1820-link-%d", rnd.Intn(MaximumClientID)))
	opts.SetOrderMatters(false)

	return &Service{
		opts: opts,

		protocols: make([]protocol.Protocol, 0),

		channel: make(chan model.Data),

		IsRun: false,
	}
}

// Channel returns a channel for consuming the received data
// Consuming this channel makes MQTT service work
func (s *Service) Channel() <-chan model.Data {
	return s.channel
}

// Register registers given protocol on MQTT service
func (s *Service) Register(p protocol.Protocol) {
	if s.IsRun {
		return // there is no way to add protocol in runtime
	}

	s.protocols = append(s.protocols, p)
}

// Protocols returns list of registered protocol's names
func (s *Service) Protocols() []string {
	names := make([]string, len(s.protocols))

	for i, p := range s.protocols {
		names[i] = p.Name()
	}

	return names
}

// Run runs application.
// This function connects MQTT client and then register its topic
func (s *Service) Run() error {
	// Create an MQTT client
	/*
		Port: 1883
		CleanSession: True
		Order: True
		KeepAlive: 30 (seconds)
		ConnectTimeout: 30 (seconds)
		MaxReconnectInterval 10 (minutes)
		AutoReconnect: True
	*/
	s.opts.SetOnConnectHandler(func(client paho.Client) {
		// subscribe to protocol topics
		for _, p := range s.protocols {
			// use $share/i1820-link in front of the MQTT topic for group subscription
			if t := s.cli.Subscribe(p.RxTopic(), 0, Handler(p, s.channel)); t.Error() != nil {
				logrus.Fatalf("mqtt subscribe error: %s", t.Error())
			}
		}
	})

	s.cli = paho.NewClient(s.opts)

	// connect to the MQTT Server
	if t := s.cli.Connect(); t.Wait() && t.Error() != nil {
		return t.Error()
	}

	s.IsRun = true

	return nil
}

// Exit closes MQTT connection then closes all channels and return from all pipeline stages
func (s *Service) Exit() {
	s.IsRun = false

	// disconnect aftere 10ms
	s.cli.Disconnect(DisconnectTimeout)
}
