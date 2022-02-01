package util_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/nogurenn/cph-wallet/util"
	"github.com/stretchr/testify/assert"
)

func Test_NewNullUUID_Valid(t *testing.T) {
	// given
	id := uuid.New()

	// when
	result := util.NewNullUUID(id)

	// then
	assert.Equal(t, id, result.UUID)
	assert.True(t, result.Valid)
}

func Test_NewNullUUID_Invalid(t *testing.T) {
	// when
	result := util.NewNullUUID(uuid.Nil)

	// then
	assert.Equal(t, uuid.Nil, result.UUID)
	assert.False(t, result.Valid)
}
