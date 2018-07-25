package actions

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aiotrc/uplink/lora"
	log "github.com/sirupsen/logrus"
)

type decodeReq struct {
	data    []byte
	project string
	device  string

	resp chan decodeRes
}

type decodeRes struct {
	data interface{}
	err  error
}

var decodeChan chan decodeReq

func init() {
	decodeChan = make(chan decodeReq, 1024)

	for i := 0; i < 5; i++ {
		go func() {
			for req := range decodeChan {
				res := decodeRes{}
				res.data, res.err = pm.RunnersDecode(req.data, req.project, req.device)
				req.resp <- res
			}
		}()
	}
}

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
		"topic":     string(topicName),
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
		"topic":     string(topicName),
	}).Info(m)

	var bdoc interface{}

	// Find thing
	p, err := pm.ThingsShow(m.DevEUI)
	if err != nil {
		log.WithFields(log.Fields{
			"component": "uplink",
		}).Errorf("PM GetThingProject: %s", err)
		return
	}

	defer func() {
		log.WithFields(log.Fields{
			"component": "uplink",
		}).Info("Insert into database")

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

	// Decode
	respChan := make(chan decodeRes)
	decodeChan <- decodeReq{
		data:    m.Data,
		project: p.Name,
		device:  m.DevEUI,

		resp: respChan,
	}
	resp := <-respChan
	close(respChan)

	bdoc = resp.data
	if resp.err != nil {
		log.WithFields(log.Fields{
			"component": "uplink",
		}).Errorf("Decode: %s", resp.err)
		return
	}
}
