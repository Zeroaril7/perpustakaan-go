package utils

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func Test_HashPassword(t *testing.T) {
	assert.NotEqual(t, HashPassword("test"), "")
}

func Test_CheckPasswordHash(t *testing.T) {
	hash := HashPassword("test")
	assert.Equal(t, CheckPasswordHash("test", hash), true)
	assert.Equal(t, CheckPasswordHash("test2", hash), false)
}
