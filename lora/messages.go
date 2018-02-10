package lora

import "time"

// RxMessage contains payloads received from your nodes
type RxMessage struct {
	ApplicationID   string
	ApplicationName string
	DeviceName      string
	DevEUI          string
	FPort           int
	FCnt            int
	RxInfo          []RxInfo
	Data            []byte
}

// RxInfo contains gateway infomation that payloads
// received from it.
type RxInfo struct {
	Mac     string
	Name    string
	Time    time.Time
	RSSI    int `json:"rssi"`
	LoRaSNR int `json:"LoRaSNR"`
}

// TxMessage contains payload send to your nodes
type TxMessage struct {
	Reference string
	Confirmed bool
	FPort     int
	Data      []byte
}
