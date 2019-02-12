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
	"fmt"
	"testing"
	"time"

	json "github.com/json-iterator/go"

	"github.com/I1820/types"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

const tID = "el-thing" // ThingID
const aName = "memory" // Asset Name
const pName = "her"    // Project Name

func TestPipelineDirect(t *testing.T) {
	a := new()
	a.run()
	ts := time.Now()

	b, err := json.Marshal(types.State{
		Raw:     18.20,
		At:      ts,
		Asset:   aName,
		ThingID: tID,
		Project: pName,
	})
	assert.NoError(t, err)
	a.stateChan.Publish(
		"i1820_fanout_states", // exchange type
		"",                    // routing key
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        b,
		},
	)
	a.Exit()

	var d types.State
	q := a.db.Collection(fmt.Sprintf("data.%s.%s", pName, tID)).FindOne(context.Background(), bson.NewDocument(
		bson.EC.SubDocument("at", bson.NewDocument(
			bson.EC.Time("$gte", ts),
		)),
		bson.EC.String("asset", aName),
	))
	assert.NoError(t, q.Decode(&d))

	assert.Equal(t, d.At.Unix(), ts.Unix())
	assert.Equal(t, 18.20, d.Value.Number)
}
