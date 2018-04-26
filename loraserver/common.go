/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 25-04-2018
 * |
 * | File Name:     loraserver/common.go
 * +===============================================
 */

package loraserver

import (
	"time"

	"github.com/brocaar/lora-app-server/api"
)

// GatewayFrame is loraserver.io raw gateway frames
type GatewayFrame struct {
	// Gateway MAC
	Mac string

	// Contains zero or one uplink frame.
	UplinkFrame *api.UplinkFrameLog

	// Contains zero or one downlink frame.
	DownlinkFrame *api.DownlinkFrameLog

	// Record creation time
	Timestamp time.Time
}
