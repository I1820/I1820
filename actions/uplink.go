package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aiotrc/uplink/decoder"
	"github.com/aiotrc/uplink/lora"
	log "github.com/sirupsen/logrus"
)

// Error handles errors
func Error(topicName, message []byte) {
	var m lora.ErrorMessage
	if err := json.Unmarshal(message, &m); err != nil {
		log.WithFields(log.Fields{
			"component": "uplink",
		}).Errorf("JSON Unmarshal: %s", err)
		return
	}

	log.WithFields(log.Fields{
		"component": "uplink",
	}).Info(m)

	if _, err := db.Collection("lora").InsertOne(context.Background(), &struct {
		Error     string
		Timestamp time.Time
		Type      string
		Project   string
		FCnt      int
	}{
		Error:     m.Error,
		Timestamp: time.Now(),
		Project:   m.ApplicationName,
		Type:      m.Type,
		FCnt:      m.FCnt,
	}); err != nil {
		log.WithFields(log.Fields{
			"component": "uplink",
		}).Errorf("Mongo insert: %s\n", err)
		return
	}

}

// Data handles uplink data
func Data(topicName, message []byte) {
	var m lora.RxMessage
	if err := json.Unmarshal(message, &m); err != nil {
		log.WithFields(log.Fields{
			"component": "uplink",
		}).Errorf("JSON Unmarshal: %s", err)
		return
	}
	log.WithFields(log.Fields{
		"component": "uplink",
	}).Info(m)

	var bdoc interface{}

	// Find thing
	p, err := pm.GetThingProject(m.DevEUI)
	if err != nil {
		log.WithFields(log.Fields{
			"component": "uplink",
		}).Errorf("PM GetThingProject: %s", err)
		return
	}

	defer func() {
		log.WithFields(log.Fields{
			"component": "uplink",
		}).Info("Insert into databse")

		if _, err := db.Collection("data").InsertOne(context.Background(), &struct {
			Raw       []byte
			Data      interface{}
			Timestamp time.Time
			ThingID   string
			RxInfo    []lora.RxInfo
			TxInfo    lora.TxInfo
			Project   string
		}{
			Raw:       m.Data,
			Data:      bdoc,
			Timestamp: time.Now(),
			ThingID:   m.DevEUI,
			RxInfo:    m.RxInfo,
			TxInfo:    m.TxInfo,
			Project:   p.Name,
		}); err != nil {
			log.WithFields(log.Fields{
				"component": "uplink",
			}).Errorf("Mongo insert: %s\n", err)
			return
		}
	}()

	// Create decoder
	decoder := decoder.New(fmt.Sprintf("http://%s:%s", Config.Decoder.Host, p.Runner.Port))

	// Decode
	parsed, err := decoder.Decode(m.Data, m.DevEUI)
	if err != nil {
		log.WithFields(log.Fields{
			"component": "uplink",
		}).Errorf("Decode: %s", err)
		return
	}

	if err := json.Unmarshal([]byte(parsed), &bdoc); err != nil {
		log.WithFields(log.Fields{
			"component": "uplink",
		}).Errorf("Unmarshal JSON: %s\n %q", err, parsed)
		return
	}

}
