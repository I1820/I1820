/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 02-08-2018
 * |
 * | File Name:     pipeline.go
 * +===============================================
 */

package core

import (
	"context"
	"fmt"
	"runtime"

	"github.com/sirupsen/logrus"
)

// insertStage inserts each data to the rabbitmq
func (a *Application) insertStage() {
	// This thread is mine
	runtime.LockOSThread()

	a.Logger.WithFields(logrus.Fields{
		"component": "dm",
	}).Info("Insert pipeline stage")

	for d := range a.insertStream {
		if _, err := a.db.Collection(fmt.Sprintf("data.%s.%s", d.Project, d.ThingID)).InsertOne(context.Background(), *d); err != nil {
			a.Logger.WithFields(logrus.Fields{
				"component": "link",
				"asset":     d.Asset,
				"thingid":   d.ThingID,
			}).Errorf("Mongo Insert: %s", err)
		} else {
			a.Logger.WithFields(logrus.Fields{
				"component": "link",
				"asset":     d.Asset,
				"thingid":   d.ThingID,
			}).Infof("Insert into database with value: %+v", d.Value)
		}
	}

	a.Logger.WithFields(logrus.Fields{
		"component": "link",
	}).Info("Insert pipeline stage is going")
	a.insertWG.Done()
}
