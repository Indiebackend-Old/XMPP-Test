package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdGen(t *testing.T) {
	id := RandomID()
	assert.Equal(t, IDLen, len(id))
}
