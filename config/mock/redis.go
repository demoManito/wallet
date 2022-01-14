package mock

import (
	"wallet/config"
	goloader "wallet/pkg/core/loader"
	"wallet/pkg/mock"
)

func CacheOption() mock.Option {
	return func(loader goloader.ILoader, me mock.IMockEnv) {
		conf, _ := config.GetMockConfig()
		rdb := config.InitRedis(conf.Redis)
		loader.Register("wallet.redis", rdb)
	}
}
