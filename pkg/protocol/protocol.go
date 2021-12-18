package protocol

import "github.com/I1820/I1820/internal/model"

// Protocol is a uplink/downlink protocol like lan or lora.
// it specifies that how we are going to receive their packets,
// for example lora packets are coming from mqtt in json.
type Protocol interface {
	TxTopic() string
	RxTopic() string

	Name() string

	Marshal([]byte) (model.Data, error)
}
