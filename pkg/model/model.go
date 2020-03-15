package model

// Model is a decoder/encoder interface like generic (based on user scripts) or aolab
type Model interface {
	Decode([]byte) interface{}
	Encode(interface{}) []byte

	Name() string
}
