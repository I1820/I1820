/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 12-05-2018
 * |
 * | File Name:     plog/plog.go
 * +===============================================
 */

package models

import "time"

// ProjectLog represents project logs
type ProjectLog struct {
	Time      time.Time `bson:"Time"`
	Message   string    `bson:"Message"`
	Code      string    `bson:"code"`
	Job       string    `bson:"job"`
	Project   string    `bson:"project"`
	Component string    `bson:"component"`
	Level     int       `bson:"Level"`
}

// LoraLog represents lora errors
type LoraLog struct {
	Timestamp time.Time
	Error     string
	Type      string
	Project   string
	FCnt      int
}
