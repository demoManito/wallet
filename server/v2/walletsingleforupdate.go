package v1

import (
	"net/http"
	"wallet/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type moneyResquest struct {
	UserID  int64   `json:"user_id"`
	Balance float64 `json:"balance"`
}

// http://127.0.0.1:8848/api/v2/wallet/draw_money
func (h *Handler) DrawMoney(c *gin.Context) {
	params := new(moneyResquest)
	if err := c.BindJSON(params); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code: 400,
			Msg:  err.Error(),
		})
		return
	}

	order, err := params.save(h.DB, models.TypeDecr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Msg:  "ok",
		Data: order,
	})
}

// http://127.0.0.1:8848/api/v2/wallet/save_money
func (h *Handler) SaveMoney(c *gin.Context) {
	params := new(moneyResquest)
	if err := c.BindJSON(params); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code: 400,
			Msg:  err.Error(),
		})
		return
	}

	order, err := params.save(h.DB, models.TypeIncr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Msg:  "ok",
		Data: order,
	})
}

func (m *moneyResquest) save(DB *gorm.DB, walletType int) (*models.Order, error) {
	order := models.NewCreateOrder(m.UserID, 0, m.Balance, models.TypeTransfer)
	err := DB.Transaction(func(tx *gorm.DB) error {
		var wallet *models.Wallet
		wallet, err := models.ForUpdateWallet(tx, m.UserID)
		if err != nil {
			return err
		}
		// 取款检查余额是否充足
		if walletType == models.TypeDecr {
			if err := wallet.IsBalanceEnough(m.Balance); err != nil {
				return err
			}
		}

		err = wallet.UpdateWallet(tx, order.UserID, wallet.Balance-m.Balance)
		if err != nil {
			logrus.Fatalf("wallet: %v", err)
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
		if err := DB.Create(order).Error; err != nil {
			logrus.Fatalf("账单保存异常: %v", err)
		}
		return nil, err
	}

	return order, nil
}
