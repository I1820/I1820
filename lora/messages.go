package lora

// TxMessage contains payloads transmitted to your nodes
type TxMessage struct {
	FPort int
	Data  []byte
}
