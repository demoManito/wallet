package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"wallet/models"
)

type transferResquest struct {
	UserID   int64   `json:"user_id"`
	ToUserID int64   `json:"to_user_id"`
	Balance  float64 `json:"balance"`
}

// http://localhost:8848/api/v1/wallet/transfer
func (h *Handler) Transfer(c *gin.Context) {
	ctx := c.Request.Context()
	params := new(transferResquest)
	if err := c.BindJSON(params); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code: 400,
			Msg:  err.Error(),
		})
		return
	}

	order, err := h.WalletService.Transfer(params.UserID, params.ToUserID, params.Balance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code: 500,
			Msg:  err.Error(),
		})
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
