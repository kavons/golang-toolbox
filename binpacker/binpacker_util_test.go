package binpacker

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zhuangsirui/binpacker"
)

func TestAddUint16Prefix(t *testing.T) {
	bytes := []byte{1, 2}
	assert.Equal(t, binpacker.AddUint16Perfix(bytes), []byte{2, 0, 1, 2}, "Prefix error.")
}

func TestAddUint32Prefix(t *testing.T) {
	bytes := []byte{1, 2}
	assert.Equal(t, binpacker.AddUint32Perfix(bytes), []byte{2, 0, 0, 0, 1, 2}, "Prefix error.")
}

func TestAddUint64Prefix(t *testing.T) {
	bytes := []byte{1, 2}
	assert.Equal(t, binpacker.AddUint64Perfix(bytes), []byte{2, 0, 0, 0, 0, 0, 0, 0, 1, 2}, "Prefix error.")
}
