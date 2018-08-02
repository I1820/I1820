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

	"github.com/aiotrc/uplink/app"
)

// Protocol implements uplink protocol for lora
type Protocol struct {
}

func Topic() []byte {
	return []byte("application/+/device/+/rx")
}

func Marshal(message []byte) (app.Data, error) {
	var m RxMessage

	if err := json.Unmarshal(message, &m); err != nil {
		return app.Data{}, err
	}

	return app.Data{
		Raw:       m.Data,
		Data:      nil,
		Timestamp: time.Now(),
		ThingID:   m.DevEUI,
		RxInfo:    m.RxInfo,
		TxInfo:    m.TxInfo,
		Project:   "",
	}, nil
}
