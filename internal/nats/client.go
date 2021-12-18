package nats

import (
	"log"

	"github.com/I1820/I1820/config"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

// ErrorHandler calls on nats errors
func ErrorHandler(conn *nats.Conn, subscription *nats.Subscription, err error) {
	logrus.WithField("component", "nats").Error(err.Error())
}

func NewClient(cfg config.NATS) *nats.EncodedConn {
	nc, err := nats.Connect(cfg.URL, nats.ErrorHandler(ErrorHandler))
	if err != nil {
		log.Fatal(err)
	}

	c, err := nats.NewEncodedConn(nc, nats.GOB_ENCODER)
	if err != nil {
		log.Fatal(err)
	}

	return c
}
