package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type transferResquest struct {
	UserID   int64   `json:"user_id"`
	ToUserID int64   `json:"to_user_id"`
	Balance  float64 `json:"balance"`
}

// http://localhost:8848/api/v2/wallet/transfer
func (h *Handler) Transfer(c *gin.Context) {
	params := new(transferResquest)
	if err := c.BindJSON(params); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code: 400,
			Msg:  err.Error(),
		})
		return
	}

	order, err := h.WalletService.ForUpdateTransfer(params.UserID, params.ToUserID, params.Balance)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code: 400,
			Msg:  err.Error(),
			Data: order,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Msg:  "ok",
		Data: order,
	})
}
