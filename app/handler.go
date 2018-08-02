/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 02-08-2018
 * |
 * | File Name:     handler.go
 * +===============================================
 */

package app

import (
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

func (a *Application) mqttHandler(p Protocol) paho.MessageHandler {
	marshaler := p.Marshal

	return func(client paho.Client, message paho.Message) {
		d, err := marshaler(message.Payload())
		if err != nil {
			a.Logger.WithFields(logrus.Fields{
				"component": "uplink",
				"topic":     message.Topic(),
			}).Errorf("Marshal error %s", err)
			return
		}
		a.Logger.WithFields(logrus.Fields{
			"component": "uplink",
			"topic":     message.Topic(),
		}).Infof("Marshal on %v", d)
		a.projectStream <- d
	}
}
