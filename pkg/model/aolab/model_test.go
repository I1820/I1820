/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 22-08-2018
 * |
 * | File Name:     model_test.go
 * +===============================================
 */

package aolab

import "testing"

const (
	rawLog = `{
		"type": "multisensor",
		"device": "hasht",
		"states": {
			"101": 101,
			"202": 202,
			"303": 303,
			"1820": 1820
		}
	}
	`
)

func TestDecode(t *testing.T) {
	m := Model{}

	r := m.Decode([]byte(rawLog))
	if r == nil {
		t.Fatalf("Invalid decode result")
	}
}

func BenchmarkDecode(b *testing.B) {
	m := Model{}

	for n := 0; n < b.N; n++ {
		m.Decode([]byte(rawLog))
	}
}
