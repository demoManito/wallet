package handler

import (
	"wallet/config"
	v1 "wallet/server/v1"
	v2 "wallet/server/v2"
)

type IHandler interface {
	GetHandler() interface{}
}

type Handler struct {
	Config *config.Config `load:"wallet.config"`

	HandlerV1 *v1.Handler `load:"wallet.handler.v1"`
	HandlerV2 *v2.Handler `load:"wallet.handler.v2"`
}

func GetHandler() *Handler {
	return new(Handler)
}
