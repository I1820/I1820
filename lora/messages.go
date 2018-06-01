package lora

// TxMessage contains payloads transmitted to your nodes
type TxMessage struct {
	Reference string // reference which will be used on ack or error (this can be a random string)
	FPort     int    // FPort to use (must be > 0)
	Data      []byte // base64 encoded data (plaintext, will be encrypted by LoRa Server)
	Confirmed bool   // whether the payload must be sent as confirmed data down or not
}
