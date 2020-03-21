package runner

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {
	m, err := New()
	assert.NoError(t, err)

	r, err := m.New(context.Background(), "Raha", nil)
	assert.NoError(t, err)

	t.Log(r.ID)

	assert.NoError(t, m.Remove(context.Background(), r))
}
