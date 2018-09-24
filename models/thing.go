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
	ID     string   `json:"id" bson:"_id,omitempty"` // thing unique identifier
	Name   string   `json:"name" bson:"name"`        // thing human readable name
	Status bool     `json:"status" bson:"status"`    // active/inactive
	Model  string   `json:"model" bson:"model"`      // model describes how to decode an incoming payload
	Tokens []string `json:"tokens" bson:"tokens"`    // thing access tokens that are generated based on K-Sortable Globally Unique IDs

	Assets         map[string]Asset       `json:"assets" bson:"assets"`                 // thing assets (sensors and actuators)
	Connectivities map[string]interface{} `json:"connectivities" bson:"connectivities"` // supported connectivities (like TheThingsNetwork and etc.)

	Project string `json:"project" bson:"project"`
}
