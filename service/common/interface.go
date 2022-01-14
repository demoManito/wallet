package common

import "wallet/models"

type IWalletService interface {
	Transfer(int64, int64, float64) (*models.Order, error)
	Serialize(interface{}) (string, error)
	ForUpdateTransfer(int64, int64, float64) (*models.Order, error)
}
