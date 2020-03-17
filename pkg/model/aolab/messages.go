/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 15-08-2018
 * |
 * | File Name:     messages.go
 * +===============================================
 */

package aolab

import "time"

// Log represents data that is coming from AoLab nodes
type Log struct {
	Timestamp time.Time
	Type      string
	Device    string
	States    map[string]interface{}
}

// Notification represents data that is going to AoLab nodes
type Notification struct {
	Type     string
	Device   string
	Settings map[string]interface{}
}
