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
	"os"
	"testing"
	"time"

	json "github.com/json-iterator/go"
	"github.com/mongodb/mongo-go-driver/bson/primitive"

	"github.com/I1820/types"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/suite"
)

// TestSuite is a test suite for core.
type TestSuite struct {
	suite.Suite

	databaseURL string
	rabbitURL   string
}

// SetupSuite initiates a test suite
func (suite *TestSuite) SetupSuite() {
	suite.rabbitURL = os.Getenv("I1820_DM_CORE_BROKER_ADDR")
	if suite.rabbitURL == "" {
		suite.rabbitURL = "amqp://admin:admin@localhost:5672/"
	}

	suite.databaseURL = os.Getenv("I1820_DM_DATABASE_URL")
	if suite.databaseURL == "" {
		suite.databaseURL = "mongodb://127.0.0.1:27017"
	}
}

// Let's test!
func TestCore(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

const tID = "el-thing" // ThingID
const aName = "memory" // Asset Name
const pName = "her"    // Project Name

func (suite *TestSuite) TestPipelineDirect() {
	a := New(suite.databaseURL, suite.rabbitURL)
	suite.NoError(a.Run())
	ts := time.Now()

	b, err := json.Marshal(types.State{
		Raw:     18.20,
		At:      ts,
		Asset:   aName,
		ThingID: tID,
		Project: pName,
	})
	suite.NoError(err)
	suite.NoError(a.ch.Publish(
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
	suite.NoError(q.Decode(&d))

	suite.Equal(d.At.Unix(), ts.Unix())
	suite.Equal(18.20, d.Raw)
}
