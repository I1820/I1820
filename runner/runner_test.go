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

import (
	"fmt"
	"testing"
)

func TestCreate(t *testing.T) {
	r := New("Eli")
	fmt.Println(r)
}
