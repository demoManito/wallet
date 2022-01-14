package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatRedisKey(t *testing.T) {
	assert := assert.New(t)

	key := FormatRedisKey("aa:", "cc")
	assert.Equal("aa:cc", key)

}
