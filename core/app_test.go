/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 03-08-2018
 * |
 * | File Name:     app_test.go
 * +===============================================
 */

package core

import (
	"context"
	"testing"
	"time"

	"github.com/I1820/types"
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// TestSuite is a test suite for core.
type TestSuite struct {
	suite.Suite

	addr string
	pm   *pm.PM
}

func TestPipeline(t *testing.T) {
	a := New()
	a.Run()
	ts := time.Now()

	a.projectStream <- types.Data{
		Raw:       []byte("Hello"),
		Data:      nil,
		Timestamp: ts,
		ThingID:   "el-thing",
		RxInfo:    nil,
		TxInfo:    nil,
		Project:   "",
	}
	time.Sleep(1 * time.Second)

	var d types.Data
	q := a.db.Collection("data").FindOne(context.Background(), bson.NewDocument(
		bson.EC.SubDocument("timestamp", bson.NewDocument(
			bson.EC.Time("$gte", ts),
		)),
		bson.EC.String("thingid", "el-thing"),
	))
	assert.NoError(t, q.Decode(&d))

	assert.Equal(t, d.Timestamp.Unix(), ts.Unix())
}

func BenchmarkPipeline(b *testing.B) {
	a := New()
	a.Run()

	wait := make(chan struct{})
	a.cli.Subscribe("i1820/project/her/raw", 0, func(client paho.Client, message paho.Message) {
		wait <- struct{}{}
	})

	for i := 0; i < b.N; i++ {
		ts := time.Now()

		a.projectStream <- types.Data{
			Raw:       []byte("Hello"),
			Data:      nil,
			Timestamp: ts,
			ThingID:   "el-thing",
			RxInfo:    nil,
			TxInfo:    nil,
			Project:   "her",
		}

		<-wait
	}
}
