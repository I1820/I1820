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

import "github.com/brocaar/lora-app-server/api"

// GatewayFrame is loraserver.io raw gateway frames
type GatewayFrame struct {
	// Gateway MAC
	Mac string

	// Contains zero or one uplink frame.
	UplinkFrames []*api.UplinkFrameLog
	// Contains zero or one downlink frame.
	DownlinkFrames []*api.DownlinkFrameLog
}
