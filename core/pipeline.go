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
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/sirupsen/logrus"
)

func (a *Application) project() {
	// this thread is mine
	runtime.LockOSThread()

	logrus.WithFields(logrus.Fields{
		"component": "link",
	}).Info("project pipeline stage has started")

	for d := range a.projectStream {
		// find the thing in I1820/pm
		t, err := a.tm.Show(d.ThingID)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"component": "link",
			}).Errorf("tm show: %s", err)
		} else {
			d.Project = t.Project
			d.Model = t.Model
		}

		if d.Project != "" && d.Model == "generic" {
			// publish raw data
			b, err := json.Marshal(d)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"component": "link",
				}).Errorf("marshal data error: %s", err)
			}
			a.cli.Publish(fmt.Sprintf("i1820/project/%s/raw", d.Project), 0, false, b)
			logrus.WithFields(logrus.Fields{
				"component": "link",
			}).Infof("publish raw data: %s", d.Project)
		}

		a.decodeStream <- d
	}
	a.projectWG.Done()
}

func (a *Application) decode() {
	// this thread is mine
	runtime.LockOSThread()

	logrus.WithFields(logrus.Fields{
		"component": "link",
	}).Info("decode pipeline stage has started")

	for d := range a.decodeStream {
		// run decode when data is coming from thing with project and it needs decode
		if d.Project != "" && d.Data == nil {
			if d.Model != "generic" {
				m, ok := a.models[d.Model]
				if !ok {
					logrus.WithFields(logrus.Fields{
						"component": "link",
					}).Errorf("model %s not found (setting the model will improves performance)", d.Model)
					// data will be parsed in project docker and pushed into mqtt parsed channel
				} else {
					d.Data = m.Decode(d.Raw)

					// publish parsed data
					b, err := json.Marshal(d)
					if err != nil {
						logrus.WithFields(logrus.Fields{
							"component": "link",
						}).Errorf("marshal data error: %s", err)
					}
					a.cli.Publish(fmt.Sprintf("i1820/project/%s/data", d.Project), 0, false, b)
					logrus.WithFields(logrus.Fields{
						"component": "link",
					}).Infof("publish parsed data: %s", d.Project)
				}
			}
		}
		a.insertStream <- d
	}
	a.decodeWG.Done()
}

func (a *Application) insert() {
	// this thread is mine
	runtime.LockOSThread()

	logrus.WithFields(logrus.Fields{
		"component": "link",
	}).Info("insert pipeline stage has started")

	for d := range a.insertStream {
		if err := a.Store.Insert(context.Background(), d); err != nil {
			logrus.WithFields(logrus.Fields{
				"component": "link",
			}).Errorf("insert into mongodb error: %s", err)
		} else {
			logrus.WithFields(logrus.Fields{
				"component": "link",
			}).Infof("insert into mongodb: %#v", d)
		}
	}
	a.insertWG.Done()
}
