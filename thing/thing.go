/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 07-02-2018
 * |
 * | File Name:     thing/thing.go
 * +===============================================
 */

package thing

// Thing contains identification and parent project
type Thing struct {
	ID     string `json:"id" bson:"id"`         // DevEUI
	Status bool   `json:"status" bson:"status"` // active/deactive
}
