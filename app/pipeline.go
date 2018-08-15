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

	"github.com/sirupsen/logrus"
)

func (a *Application) project() {
	a.Logger.WithFields(logrus.Fields{
		"component": "uplink",
	}).Info("Project pipeline stage")

	for d := range a.projectStream {
		// Find thing project in I1820/pm
		p, err := a.pm.ThingsShow(d.ThingID)
		if err != nil {
			a.Logger.WithFields(logrus.Fields{
				"component": "uplink",
			}).Errorf("PM ThingsShow: %s", err)
		} else {
			d.Project = p.Name
			d.Model = p.Model
		}
		a.decodeStream <- d
	}
}

func (a *Application) decode() {
	a.Logger.WithFields(logrus.Fields{
		"component": "uplink",
	}).Info("Decode pipeline stage")

	for d := range a.decodeStream {
		// Run decode when data is comming from thing with project and its needs decode
		if d.Project != "" && d.Data == nil {
			if d.Model == "generic" {
				data, err := a.pm.RunnersDecode(d.Raw, d.Project, d.ThingID)
				if err != nil {
					a.Logger.WithFields(logrus.Fields{
						"component": "uplink",
					}).Errorf("PM RunnersDecode: %s", err)
				} else {
					d.Data = data
				}
			} else {
				m, ok := a.models[d.Model]
				if !ok {
					a.Logger.WithFields(logrus.Fields{
						"component": "uplink",
					}).Errorf("Model %s not found", d.Model)
				} else {
					d.Data = m.Decode(d.Raw)
				}
			}
		}
		a.insertStream <- d
	}
}

func (a *Application) insert() {
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
