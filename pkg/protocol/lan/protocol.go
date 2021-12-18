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

package lan

import (
	"encoding/json"
	"time"

	"github.com/I1820/I1820/internal/model"
	lan "github.com/I1820/lanserver/model"
)

// Protocol implements uplink protocol for lora
type Protocol struct {
}

// RxTopic returns lan rx message topic
func (p Protocol) RxTopic() string {
	return "device/+/rx"
}

// TxTopic returns lan tx message topic
func (p Protocol) TxTopic() string {
	return "device/+/tx"
}

// Name returns protocol unique name
func (p Protocol) Name() string {
	return "lan"
}

// Marshal marshals given lan byte message (in json format) into platform data structure
func (p Protocol) Marshal(message []byte) (model.Data, error) {
	var m lan.RxMessage

	if err := json.Unmarshal(message, &m); err != nil {
		return model.Data{}, err
	}

	return model.Data{
		Raw:       m.Data,
		Data:      nil,
		Timestamp: time.Now(),
		ThingID:   m.DevEUI,
		RxInfo:    nil,
		TxInfo:    nil,
		Project:   "",
	}, nil
}
