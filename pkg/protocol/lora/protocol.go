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

	"github.com/I1820/I1820/model"
)

// Protocol implements uplink protocol for lora
type Protocol struct {
}

// RxTopic returns lora rx message topic
// https://www.loraserver.io/lora-app-server/integrate/sending-receiving/mqtt/
func (p Protocol) RxTopic() string {
	return "application/+/device/+/rx"
}

// TxTopic returns lora tx message topic for each downlink message
// https://www.loraserver.io/lora-app-server/integrate/sending-receiving/mqtt/
func (p Protocol) TxTopic() string {
	return "application/+/device/+/tx"
}

// Name returns protocol unique name
func (p Protocol) Name() string {
	return "lora"
}

// Marshal marshals given lora byte message (in json format) into platform data structure
// https://www.loraserver.io/lora-app-server/integrate/sending-receiving/mqtt/
func (p Protocol) Marshal(message []byte) (model.Data, error) {
	var m RxMessage

	if err := json.Unmarshal(message, &m); err != nil {
		return model.Data{}, err
	}

	return model.Data{
		Raw:       m.Data,
		Data:      nil,
		Timestamp: time.Now(),
		ThingID:   m.DevEUI,
		RxInfo:    m.RxInfo,
		TxInfo:    m.TxInfo,
		Project:   "",
	}, nil
}
