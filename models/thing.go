package models

// Thing contains thing identification and its parent project
type Thing struct {
	Name    string `json:"name" bson:"name"`       // DevEUI
	Status  bool   `json:"status" bson:"status"`   // active/inactive
	Model   string `json:"model" bson:"model"`     // model describes project decoder
	Project string `json:"project" bson:"project"` // project identification
}
