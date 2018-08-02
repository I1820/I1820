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
	"github.com/sirupsen/logrus"
)

func (a *Application) mqttHandler(p Protocol) func(topicName, message []byte) {
	marshaler := p.Marshal

	return func(topicName, message []byte) {
		d, err := marshaler(message)
		if err != nil {
			a.Logger.WithFields(logrus.Fields{
				"component": "uplink",
				"topic":     string(topicName),
			}).Errorf("Marshal error %s", err)
			return
		}
		a.Logger.WithFields(logrus.Fields{
			"component": "uplink",
			"topic":     string(topicName),
		}).Infof("Marshal on %v", d)
		a.projectStream <- d
	}
}
