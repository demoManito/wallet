package work

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"wallet/models"
	goloader "wallet/pkg/core/loader"
)

// LargeVolume 大交易量钱包
var LargeVolume = []string{
	"", "111", "222", "333", "444", "555", "666",
}

func init() {
	goloader.Register(func(loader goloader.ILoader) error {
		walletManger := new(WalletMangerWorker)
		return loader.Register("wallet.worker.WalletMangerWorker", walletManger)
	})
}

type WalletMangerWorker struct {
	DB    *gorm.DB        `load:"wallet.db"`
	Redis *redis.Client   `load:"wallet.redis"`
	Ctx   context.Context `load:"wallet.Context"`
}

// WalletMangerSubscribe 遍历转账大户
// TODO: LargeVolume 接口入队需要重新处理下
func (w *WalletMangerWorker) WalletMangerSubscribe() {
	for _, id := range LargeVolume {
		pubsub := w.Redis.Subscribe(w.Ctx, models.FormatWalletManage(id))
		go func() {
			for {
				select {
				case data := <-pubsub.Channel():
					// 金额扣减增加
					logrus.Infof("chan: %s; data: %s", data.Channel, data.Payload)

					order := new(models.Order)
					if err := json.Unmarshal([]byte(data.Payload), order); err != nil {
						continue
					}
					err := w.consume(order)
					if err != nil {
						logrus.Fatalf("[WalletMangerSubscribe] 账单处理异常: %v", err)
					}
				case <-w.Ctx.Done():
					logrus.Error("[WalletMangerSubscribe] Done")
					return
				}
			}
		}()
	}
}

func (w *WalletMangerWorker) consume(order *models.Order) (err error) {
	switch order.Type {
	case models.TypeIncr:
		err = w.recharge(order)
	case models.TypeDecr:
		err = w.drawMoney(order)
	case models.TypeTransfer:
		err = w.transfer(order)
	}
	return
}

func (w *WalletMangerWorker) transfer(order *models.Order) error {
	err := w.DB.Transaction(func(tx *gorm.DB) error {
		balance := order.Balance

		// 支付钱包
		payWallet, err := models.FindWalletByUserID(tx, order.UserID)
		if err != nil {
			return err
		}
		if payWallet.Balance-balance < 0 {
			return models.InsufficientBalanceErr
		}
		// 收款钱包
		collectionWallet, err := models.FindWalletByUserID(tx, order.ToUserID)
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

		err = order.SaveOrderStatus(tx, order.ID, models.StatusSuccess)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		logrus.Fatalf("转账异常: %v", err)
		if err := order.SaveOrderStatus(w.DB, order.ID, models.StatusErr); err != nil {
			logrus.Fatalf("账单保存异常: %v", err)
		}
		return err
	}

	return nil
}

func (w *WalletMangerWorker) drawMoney(order *models.Order) error {
	err := w.DB.Transaction(func(tx *gorm.DB) error {
		balance := order.Balance

		wallet, err := models.FindWalletByUserID(tx, order.UserID)
		if err != nil {
			return err
		}
		// 余额不足
		if err := wallet.IsBalanceEnough(balance); err != nil {
			return err
		}

		err = wallet.UpdateWallet(tx, order.UserID, wallet.Balance-balance)
		if err != nil {
			return err
		}

		err = order.SaveOrderStatus(tx, order.ID, models.StatusSuccess)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		logrus.Fatalf("取款异常: %v", err)
		if err = order.SaveOrderStatus(w.DB, order.ID, models.StatusErr); err != nil {
			logrus.Fatalf("账单保存异常: %v", err)
		}
		return err
	}

	return nil
}

// Recharge
func (w *WalletMangerWorker) recharge(order *models.Order) error {
	err := w.DB.Transaction(func(tx *gorm.DB) error {
		wallet, err := models.FindWalletByUserID(tx, order.UserID)
		if err != nil {
			return err
		}
		err = wallet.UpdateWallet(tx, order.UserID, wallet.Balance+order.Balance)
		if err != nil {
			return err
		}

		err = order.SaveOrderStatus(tx, order.ID, models.StatusSuccess)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		logrus.Fatalf("存款异常: %v", err)
		if err = order.SaveOrderStatus(w.DB, order.ID, models.StatusErr); err != nil {
			logrus.Fatalf("账单保存异常: %v", err)
		}
		return err
	}

	return nil
}
