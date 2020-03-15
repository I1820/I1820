/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 13-08-2018
 * |
 * | File Name:     data.go
 * +===============================================
 */

package model

import "time"

// Data represents nodes data after parse phase that is done by protocols.
// Each data have its raw and decoded payload and may have link quality information.
// This structure is created and remains in the platform for each incoming data.
type Data struct {
	Raw       []byte      `json:"raw" bson:"raw"`             // data before decode
	Data      interface{} `json:"data" bson:"data"`           // data after decode
	Timestamp time.Time   `json:"timestamp" bson:"timestamp"` // when data received in uplink
	ThingID   string      `json:"thing_id" bson:"thingid"`    // deveui

	RxInfo interface{} `json:"rx_info" bson:"rxinfo"`
	TxInfo interface{} `json:"tx_info" bson:"txinfo"`

	Project  string `json:"project" bson:"project"`   // thing project identification
	Protocol string `json:"protocol" bson:"protocol"` // uplink protocol
	Model    string `json:"model" bson:"model"`       // way of decode
}
