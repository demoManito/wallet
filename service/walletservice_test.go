package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"wallet/models"
)

func TestSerialize(t *testing.T) {
	assert := assert.New(t)

	service := &walletService{}
	_, err := service.Serialize(&models.Order{})
	assert.NoError(err)
}
