package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMax(t *testing.T) {
	var a, b = 1, 2
	assert.EqualValues(t, 2, Max(a, b))

	var strA, strB = "a", "b"
	assert.EqualValues(t, "b", Max(strA, strB))
}

func TestMin(t *testing.T) {
	var a, b = 1, 2
	assert.EqualValues(t, 1, Min(a, b))

	var strA, strB = "a", "b"
	assert.EqualValues(t, "a", Min(strA, strB))
}

func TestStripContent(t *testing.T) {
	var content = "1234567890"
	assert.EqualValues(t, "123", StripContent(content, 3))
	assert.EqualValues(t, "1234567890", StripContent(content, 100))

	content = "你好，世界"
	assert.EqualValues(t, "你好", StripContent(content, 2))
	assert.EqualValues(t, "你好，世界", StripContent(content, 100))
}
