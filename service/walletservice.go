package service

import (
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"wallet/models"
	goloader "wallet/pkg/core/loader"
	"wallet/service/common"
)

func init() {
	goloader.Register(func(loader goloader.ILoader) error {
		var service common.IWalletService = new(walletService)
		return loader.Register("wallet.service.walletService", service)
	})
}

var _ common.IWalletService = new(walletService)

type walletService struct {
	DB    *gorm.DB      `load:"wallet.db"`
	Redis *redis.Client `load:"wallet.redis"`
}

func (w *walletService) Transfer(userID, toUserID int64, balance float64) (*models.Order, error) {
	order, err := models.Transfer(w.DB, userID, toUserID, balance)
	if err != nil {
		return nil, err
	}
	return order, err
}

func (w *walletService) ForUpdateTransfer(userID, toUserID int64, balance float64) (*models.Order, error) {
	order := models.NewCreateOrder(userID, toUserID, balance, models.TypeTransfer)
	err := w.DB.Transaction(func(tx *gorm.DB) error {
		var payWallet *models.Wallet
		var collectionWallet *models.Wallet
		var err error
		if userID < toUserID {
			payWallet, err = models.ForUpdateWallet(tx, userID)
			collectionWallet, err = models.ForUpdateWallet(tx, toUserID)
		} else {
			collectionWallet, err = models.ForUpdateWallet(tx, toUserID)
			payWallet, err = models.ForUpdateWallet(tx, userID)
		}
		if err != nil {
			return err
		}
		// 余额不足
		if err := payWallet.IsBalanceEnough(balance); err != nil {
			return err
		}

		err = payWallet.UpdateWallet(tx, order.UserID, payWallet.Balance-balance)
		if err != nil {
			logrus.Fatalf("payWallet: %v", err)
			return err
		}
		err = collectionWallet.UpdateWallet(tx, order.ToUserID, collectionWallet.Balance+balance)
		if err != nil {
			logrus.Fatalf("collectionWallet: %v", err)
			return err
		}

		order.Status = models.StatusSuccess
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		order.Status = models.StatusErr
		if err := w.DB.Create(order).Error; err != nil {
			logrus.Fatalf("账单保存异常: %v", err)
		}
		return nil, err
	}

	return order, nil
}

func (w *walletService) Serialize(option interface{}) (string, error) {
	serialize, err := json.Marshal(option)
	if err != nil {
		return "", err
	}

	return string(serialize), nil
}
