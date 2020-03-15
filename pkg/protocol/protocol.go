package protocol

import "github.com/I1820/I1820/model"

// Protocol is a uplink/downlink protocol like lan or lora
type Protocol interface {
	TxTopic() string
	RxTopic() string

	Name() string

	Marshal([]byte) (model.Data, error)
}
