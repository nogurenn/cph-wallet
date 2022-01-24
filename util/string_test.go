package util_test

import (
	"testing"

	"github.com/nogurenn/cph-wallet/util"
	"github.com/stretchr/testify/assert"
)

func Test_AtLeastOneEmptyString_True(t *testing.T) {
	// when
	result := util.AtLeastOneEmptyString("hello", "", "world")

	// then
	assert.True(t, result)
}

func Test_AtLeastOneEmptyString_False(t *testing.T) {
	// when
	result := util.AtLeastOneEmptyString("hello", "world")

	// then
	assert.False(t, result)
}
