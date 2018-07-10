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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {
	r, err := New("Eli", nil)

	assert.NoError(t, err)

	t.Log(r.ID)

	assert.NoError(t, r.Remove())
}
