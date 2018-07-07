/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 07-02-2018
 * |
 * | File Name:     thing/thing.go
 * +===============================================
 */

package models

// Thing contains identification and parent project of a thing
type Thing struct {
	ID     string `json:"id" bson:"id"`         // DevEUI
	Status bool   `json:"status" bson:"status"` // active/inactive
}
