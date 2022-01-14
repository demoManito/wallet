package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveOrderStatus(t *testing.T) {
	assert := assert.New(t)

	order := NewCreateOrder(1, 0, 100, TypeTransfer)
	assert.Equal(StatusCreate, order.Status)
	assert.Equal(TypeTransfer, order.Type)

	err := order.SaveOrderStatus(tester.DB, 2, StatusSuccess)
	assert.NoError(err)

	order = new(Order)
	err = tester.DB.Where(2).Find(order).Error
	assert.NoError(err)
	assert.Equal(StatusSuccess, order.Status)
}
