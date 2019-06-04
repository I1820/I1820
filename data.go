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

package types

import "time"

// Data represents nodes data after parse phase that is done by protocols.
// Each data have its raw and decoded payload and may have link quality information.
// This structure is created and remains in the platform for each incoming data.
type Data struct {
	Raw       []byte      `json:"raw"`       // data before decode
	Data      interface{} `json:"data"`      // data after decode
	Timestamp time.Time   `json:"timestamp"` // when data received in uplink
	ThingID   string      `json:"thing_id"`  // deveui

	RxInfo interface{} `json:"rx_info"`
	TxInfo interface{} `json:"tx_info"`

	Project  string `json:"project"`  // thing project identification
	Protocol string `json:"protocol"` // uplink protocol
	Model    string `json:"model"`    // way of decode
}
