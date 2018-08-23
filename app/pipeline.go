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

package app

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/sirupsen/logrus"
)

func (a *Application) project() {
	// This thread is mine
	runtime.LockOSThread()

	a.Logger.WithFields(logrus.Fields{
		"component": "uplink",
	}).Info("Project pipeline stage")

	for d := range a.projectStream {
		// Find thing in I1820/pm
		t, err := a.pm.ThingsShow(d.ThingID)
		if err != nil {
			a.Logger.WithFields(logrus.Fields{
				"component": "uplink",
			}).Errorf("PM ThingsShow: %s", err)
		} else {
			d.Project = t.Project
			d.Model = t.Model
		}

		if d.Project != "" {
			// Publish raw data
			b, err := json.Marshal(d)
			if err != nil {
				a.Logger.WithFields(logrus.Fields{
					"component": "uplink",
				}).Errorf("Marshal data error: %s", err)
			}
			a.cli.Publish(fmt.Sprintf("i1820/project/%s/raw", d.Project), 0, false, b)
			a.Logger.WithFields(logrus.Fields{
				"component": "uplink",
			}).Infof("Publish data into runner %s", d.Project)

		}

		a.decodeStream <- d
	}
}

func (a *Application) decode() {
	// This thread is mine
	runtime.LockOSThread()

	a.Logger.WithFields(logrus.Fields{
		"component": "uplink",
	}).Info("Decode pipeline stage")

	for d := range a.decodeStream {
		// Run decode when data is comming from thing with project and it needs decode
		if d.Project != "" && d.Data == nil {
			if d.Model != "generic" {
				m, ok := a.models[d.Model]
				if !ok {
					a.Logger.WithFields(logrus.Fields{
						"component": "uplink",
					}).Errorf("Model %s not found", d.Model)
				} else {
					d.Data = m.Decode(d.Raw)
				}
				// Publish parsed data
				b, err := json.Marshal(d)
				if err != nil {
					a.Logger.WithFields(logrus.Fields{
						"component": "uplink",
					}).Errorf("Marshal data error: %s", err)
				}
				a.cli.Publish(fmt.Sprintf("i1820/project/%s/data", d.Project), 0, false, b)
				a.Logger.WithFields(logrus.Fields{
					"component": "uplink",
				}).Infof("Publish data into runner %s", d.Project)
			}
		}
		a.insertStream <- d
	}
}

func (a *Application) insert() {
	// This thread is mine
	runtime.LockOSThread()

	a.Logger.WithFields(logrus.Fields{
		"component": "uplink",
	}).Info("Insert pipeline stage")

	for d := range a.insertStream {
		if _, err := a.db.Collection("data").InsertOne(context.Background(), d); err != nil {
			a.Logger.WithFields(logrus.Fields{
				"component": "uplink",
			}).Errorf("Mongo Insert: %s", err)
		} else {
			a.Logger.WithFields(logrus.Fields{
				"component": "uplink",
			}).Infof("Insert into database: %#v", d)
		}
	}
}
