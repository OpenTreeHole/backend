package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringToBytes(t *testing.T) {
	s := "123"
	assert.EqualValues(t, []byte{'1', '2', '3'}, StringToBytes(s))
}

func TestBytesToString(t *testing.T) {
	b := []byte{'1', '2', '3'}
	assert.EqualValues(t, "123", BytesToString(b))
}
