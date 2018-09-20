/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 19-09-2018
 * |
 * | File Name:     asset.go
 * +===============================================
 */

package models

// Asset is sensor or actuator that is attached to a thing.
type Asset struct {
	Title string `json:"title" bson:"title"`
	Type  string `json:"type" bson:"type"` // boolean, number, string, object and array
	Kind  string `json:"kind" bson:"kind"` // sensor or actuator
}
