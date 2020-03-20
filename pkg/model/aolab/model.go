/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 15-08-2018
 * |
 * | File Name:     model.go
 * +===============================================
 */

package aolab

import (
	"encoding/json"
)

// Model reperesents AoLab model. this model for marshaling
// and unmarshaling of data is created originally by
// Amirkabir University IoT Lab
type Model struct{}

// Name returns model name
func (m Model) Name() string {
	return "aolab"
}

// Decode given data with aolab structure
func (m Model) Decode(d []byte) interface{} {
	var l Log

	if err := json.Unmarshal(d, &l); err != nil {
		return nil
	}

	return l.States
}

// Encode given object with aolab structure
func (m Model) Encode(o interface{}) []byte {
	n, ok := o.(Notification)
	if !ok {
		return nil
	}

	b, err := json.Marshal(n)
	if err != nil {
		return nil
	}

	return b
}
