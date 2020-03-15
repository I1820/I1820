package protocol

import "github.com/I1820/types"

// Protocol is a uplink/downlink protocol like lan or lora
type Protocol interface {
	TxTopic() string
	RxTopic() string

	Name() string

	Marshal([]byte) (types.Data, error)
}
