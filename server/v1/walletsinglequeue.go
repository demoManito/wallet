package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"wallet/models"
)

type moneyResquest struct {
	UserID  int64   `json:"user_id"`
	Balance float64 `json:"balance"`
}

// http://127.0.0.1:8848/api/v1/wallet/draw_money
func (h *Handler) DrawMoney(c *gin.Context) {
	ctx := c.Request.Context()
	params := new(moneyResquest)
	if err := c.BindJSON(params); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code: 400,
			Msg:  err.Error(),
		})
		return
	}

	order, err := params.CreateOrder(c, h.DB, params.UserID, params.Balance, models.TypeDecr)
	if err != nil {
		return
	}

	serialize, err := h.WalletService.Serialize(order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}

	err = h.Redis.Publish(ctx, models.FormatWalletManage(models.FormatWalletUserID(params.UserID)), serialize).Err()
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

// http://127.0.0.1:8848/api/v1/wallet/save_money
func (h *Handler) SaveMoney(c *gin.Context) {
	ctx := c.Request.Context()
	params := new(moneyResquest)
	if err := c.BindJSON(params); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code: 400,
			Msg:  err.Error(),
		})
		return
	}

	order, err := params.CreateOrder(c, h.DB, params.UserID, params.Balance, models.TypeIncr)
	if err != nil {
		return
	}

	serialize, err := h.WalletService.Serialize(order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}

	err = h.Redis.Publish(ctx, models.FormatWalletManage(models.FormatWalletUserID(params.UserID)), serialize).Err()
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

// CreateOrder
func (m *moneyResquest) CreateOrder(c *gin.Context, DB *gorm.DB, userID int64, balance float64, optionType int) (*models.Order, error) {
	var err error
	order := new(models.Order)

	switch optionType {
	case models.TypeIncr:
		order, err = models.SaveMoney(DB, userID, balance)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Code: 500,
				Msg:  err.Error(),
			})
			return nil, err
		}
	case models.TypeDecr:
		order, err = models.DrawMoney(DB, userID, balance)
		if err != nil {
			if err == models.InsufficientBalanceErr {
				c.JSON(http.StatusForbidden, Response{
					Code: 403,
					Msg:  err.Error(),
				})
				return nil, err
			}

			c.JSON(http.StatusInternalServerError, Response{
				Code: 500,
				Msg:  err.Error(),
			})
			return nil, err
		}
	default:
		return nil, errors.New("参数异常")
	}

	return order, nil
}
