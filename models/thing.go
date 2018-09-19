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
	Model  string `json:"model" bson:"model"`   // model describes project decoder

	Project string `json:"project,omitempty" bson:"project,omitempty"`
	User    string `json:"user,omitempty" bson:"user,omitempty"`
}
