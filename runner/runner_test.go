/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 06-02-2018
 * |
 * | File Name:     runner/runner_test.go
 * +===============================================
 */

package runner

import "testing"

func TestBasic(t *testing.T) {
	r := New("Eli")

	t.Log(r.ID)

	r.Remove()
}
