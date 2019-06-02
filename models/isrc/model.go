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

package isrc

// Model reperesents ISRC model. this model for marshaling
// and unmarshaling of data is created originally by
// Iranian Space Research Center
type Model struct{}

// Name returns model name
func (m Model) Name() string {
	return "isrc"
}

// Decode given data with aolab structure
func (m Model) Decode(d []byte) interface{} {
	return nil
}

// Encode given object with aolab structure
func (m Model) Encode(o interface{}) []byte {
	return nil
}
