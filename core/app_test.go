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
	"github.com/mongodb/mongo-go-driver/bson/primitive"

	"github.com/I1820/dm/config"
	"github.com/I1820/types"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

const tID = "el-thing" // ThingID
const aName = "memory" // Asset Name
const pName = "her"    // Project Name

func TestPipelineDirect(t *testing.T) {
	a := New(config.GetConfig().Database.URL, fmt.Sprintf("amqp://%s:%s@%s/", config.GetConfig().Core.Broker.User, config.GetConfig().Core.Broker.Pass, config.GetConfig().Core.Broker.Host))
	assert.NoError(t, a.Run())
	ts := time.Now()

	b, err := json.Marshal(types.State{
		Raw:     18.20,
		At:      ts,
		Asset:   aName,
		ThingID: tID,
		Project: pName,
	})
	assert.NoError(t, err)
	assert.NoError(t, a.stateChan.Publish(
		"i1820_fanout_states", // exchange type
		"",                    // routing key
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        b,
		},
	))
	time.Sleep(1 * time.Second) // wait for rabbitmq to publish the message.
	a.Exit()

	var d types.State
	q := a.db.Collection(fmt.Sprintf("data.%s.%s", pName, tID)).FindOne(context.Background(), primitive.M{
		"at": primitive.M{
			"$gte": ts,
		},
		"asset": aName,
	})
	assert.NoError(t, q.Decode(&d))

	assert.Equal(t, d.At.Unix(), ts.Unix())
	assert.Equal(t, 18.20, d.Raw)
}
