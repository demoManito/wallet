package models

import (
	"testing"

	"gorm.io/gorm"

	mock2 "wallet/config/mock"
	goloader "wallet/pkg/core/loader"
	"wallet/pkg/mock"
)

var tester = new(modeltest)

type modeltest struct {
	DB *gorm.DB `load:"wallet.db"`
}

func TestMain(m *testing.M) {
	loader := goloader.NewSingleLoader()
	e := mock.NewMockEnv(nil, loader, mock2.DBOption()) // mock2.CacheOption()
	defer e.Close()
	goloader.LoadingAll(loader)
	loader.Loading(tester)
	m.Run()
}
