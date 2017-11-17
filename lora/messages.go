package lora

// RxMessage contains payloads received from your nodes
type RxMessage struct {
	ApplicationID   string
	ApplicationName string
	DeviceName      string
	DevEUI          string
	FPort           int
	FCnt            int
	Data            []byte
}

// TxMessage contains payload send to your nodes
type TxMessage struct {
	Reference string
	Confirmed bool
	FPort     int
	Data      []byte
}
