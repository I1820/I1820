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

// Data represents uplink data and metadata
// this structure create for each uplink data and remain in platform
type Data struct {
	Raw       []byte      // data before decode
	Data      interface{} // data after decode
	Timestamp time.Time   // when data received in uplink
	ThingID   string      // deveui
	RxInfo    interface{}
	TxInfo    interface{}
	Project   string // thing project identification
	Protocol  string // uplink protocol
}
