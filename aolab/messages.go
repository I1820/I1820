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

// Log represents data that is comming from aolab nodes
type Log struct {
	// TODO
	//	Timestamp time.Time
	Type   string
	Device string
	States map[string]interface{}
}
