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

package app

import (
	"context"
	"testing"
	"time"

	"github.com/I1820/types"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/stretchr/testify/assert"
)

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
