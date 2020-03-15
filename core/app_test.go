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
	"os"
	"testing"
	"time"

	"github.com/I1820/types"
	paho "github.com/eclipse/paho.mqtt.golang"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
)

// TestSuite is a test suite for core.
type TestSuite struct {
	suite.Suite

	dbURL     string
	pmURL     string
	brokerURL string
}

// SetupSuite initiates a test suite
func (suite *TestSuite) SetupSuite() {
	suite.brokerURL = os.Getenv("I1820_LINK_CORE_BROKER_ADDR")
	if suite.brokerURL == "" {
		suite.brokerURL = "tcp://127.0.0.1:1883"
	}

	// on the CI tests we assume that this service is not available
	suite.pmURL = os.Getenv("I1820_LINK_PM_URL")
	if suite.pmURL == "" {
		suite.pmURL = "http://127.0.0.1:8080"
	}

	suite.dbURL = os.Getenv("I1820_LINK_DATABASE_URL")
	if suite.dbURL == "" {
		suite.dbURL = "mongodb://127.0.0.1:27017"
	}
}

// Let's test!
func TestCore(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

const tID = "el-thing" // ThingID

func (suite *TestSuite) TestPipeline() {
	a, err := New(suite.pmURL, suite.dbURL, suite.brokerURL)
	suite.Require().NoError(err)
	suite.Require().NoError(a.Run())
	ts := time.Now()

	a.projectStream <- types.Data{
		Raw:       []byte("Hello"),
		Data:      nil,
		Timestamp: ts,
		ThingID:   tID,
		RxInfo:    nil,
		TxInfo:    nil,
		Project:   "",
	}
	time.Sleep(1 * time.Second)

	var d types.Data
	q := a.db.Collection("data").FindOne(context.Background(), bson.M{
		"timestamp": bson.M{
			"$gte": ts,
		},
		"thingid": "el-thing",
	})
	suite.NoError(q.Decode(&d))

	suite.Equal(d.Timestamp.Unix(), ts.Unix())
}

func BenchmarkPipeline(b *testing.B) {
	brokerURL := os.Getenv("I1820_LINK_CORE_BROKER_ADDR")
	if brokerURL == "" {
		brokerURL = "tcp://127.0.0.1:1883"
	}

	// on the CI tests we assume that this service is not available
	pmURL := os.Getenv("I1820_LINK_PM_URL")
	if pmURL == "" {
		pmURL = "http://127.0.0.1:8080"
	}

	dbURL := os.Getenv("I1820_LINK_DATABASE_URL")
	if dbURL == "" {
		dbURL = "mongodb://127.0.0.1:27017"
	}

	a, err := New(pmURL, dbURL, brokerURL)
	require.NoError(b, err)
	require.NoError(b, a.Run())

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
