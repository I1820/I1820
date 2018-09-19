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
	ID     string `json:"id" bson:"_id,omitempty"` // Thing unique identifier
	Name   string `json:"name" bson:"name"`        // Thing human readable name
	Status bool   `json:"status" bson:"status"`    // active/inactive
	Model  string `json:"model" bson:"model"`      // model describes how to decode an incoming payload

	Project string `json:"project" bson:"project`
}
