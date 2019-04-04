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

	log "github.com/sirupsen/logrus"
)

// insertStage inserts each data to the rabbitmq
func (a *Application) insertStage() {
	// This thread is mine
	runtime.LockOSThread()

	log.WithFields(log.Fields{
		"component": "dm",
	}).Info("Insert pipeline stage")

	for d := range a.insertStream {
		if _, err :=
			a.db.Collection(fmt.Sprintf("data.%s.%s", d.Project, d.ThingID)).InsertOne(context.Background(), *d); err != nil {
			log.WithFields(log.Fields{
				"component": "dm",
				"asset":     d.Asset,
				"thingid":   d.ThingID,
			}).Errorf("Insert into database failed: %s", err)
		} else {
			log.WithFields(log.Fields{
				"component": "dm",
				"asset":     d.Asset,
				"thingid":   d.ThingID,
			}).Infof("Insert into database with value: %+v", d.Value)
		}
	}

	log.WithFields(log.Fields{
		"component": "dm",
	}).Info("Insert pipeline stage is going")
	a.insertWG.Done()
}
