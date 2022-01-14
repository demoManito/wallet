package v1

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"wallet/config"
	goloader "wallet/pkg/core/loader"
	"wallet/service/common"
)

func init() {
	goloader.Register(func(loader goloader.ILoader) error {
		handler := new(Handler)
		return loader.Register("wallet.handler.v1", handler)
	})
}

type Handler struct {
	WalletService common.IWalletService `load:"wallet.service.walletService"`

	DB     *gorm.DB       `load:"wallet.db"`
	Redis  *redis.Client  `load:"wallet.redis"`
	Config *config.Config `load:"wallet.config"`
}

// Response
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func (*Handler) GetHandler() interface{} {
	return new(Handler)
}
