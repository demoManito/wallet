package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestFindWalletByUserID(t *testing.T) {
	assert := assert.New(t)

	wallet, err := FindWalletByUserID(tester.DB, 1)
	assert.NoError(err)
	assert.NotNil(wallet)

	wallet, err = FindWalletByUserID(tester.DB, 10000000)
	assert.Error(err, gorm.ErrRecordNotFound.Error())
	assert.Nil(wallet)
}

func TestFormatWalletManage(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("wallet/manage/1", FormatWalletManage(1))

	assert.Equal("wallet/manage", FormatWalletManage(""))
}
