package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"wallet/models"
)

type statusRequest struct {
	OrderID int64 `json:"order_id"`
	UserID  int64 `json:"user_id"`
}

func (h *Handler) Status(c *gin.Context) {
	params := new(statusRequest)
	if err := c.BindJSON(params); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code: 400,
			Msg:  err.Error(),
		})
		return
	}

	order := &models.Order{
		Model: gorm.Model{
			ID: uint(params.OrderID),
		},
		UserID: params.UserID,
	}
	err := h.DB.Where(order).First(order).Error
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
