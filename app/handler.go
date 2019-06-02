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
			logrus.WithFields(logrus.Fields{
				"component": "link",
				"topic":     message.Topic(),
			}).Errorf("marshal error %s", err)
			return
		}
		d.Protocol = p.Name()
		logrus.WithFields(logrus.Fields{
			"component": "link",
			"topic":     message.Topic(),
		}).Infof("marshal on %v", d)
		a.projectStream <- d
	}
}
