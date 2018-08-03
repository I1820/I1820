/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 02-08-2018
 * |
 * | File Name:     protocol.go
 * +===============================================
 */

package lora

import (
	"encoding/json"
	"time"

	"github.com/aiotrc/lanserver/models"
	"github.com/aiotrc/uplink/app"
)

// Protocol implements uplink protocol for lora
type Protocol struct {
}

// Topic returns lan message topic
func (p Protocol) Topic() string {
	return "device/+/rx"
}

// Marshal marshals given lan byte message (in json format) into platform data structure
func (p Protocol) Marshal(message []byte) (app.Data, error) {
	var m models.RxMessage

	if err := json.Unmarshal(message, &m); err != nil {
		return app.Data{}, err
	}

	return app.Data{
		Raw:       m.Data,
		Data:      nil,
		Timestamp: time.Now(),
		ThingID:   m.DevEUI,
		RxInfo:    nil,
		TxInfo:    nil,
		Project:   "",
	}, nil
}
