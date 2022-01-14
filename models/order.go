package models

import "gorm.io/gorm"

const (
	TypeIncr     = iota + 1 // 1 充值
	TypeDecr                // 2 取款
	TypeTransfer            // 3 转账
)

const (
	StatusCreate  = iota // 订单创建
	StatusSuccess        // 订单成功
	StatusCancel         // 订单取消
	StatusErr            // 订单失败
)

type Order struct {
	gorm.Model

	UserID   int64   `gorm:"column:user_id" json:"user_id"`
	ToUserID int64   `gorm:"column:to_user_id" json:"to_user_id"` // 转账用户ID
	Balance  float64 `gorm:"column:balance" json:"balance"`       // 转账金额
	Type     int     `gorm:"column:type" json:"type"`             // 操作类型
	Status   int     `gorm:"column:status" json:"status"`         // 订单状态
	Msg      string  `gorm:"column:msg" json:"msg"`               // 订单详情
}

func (*Order) TableName() string {
	return "orders"
}

func NewCreateOrder(userID, toUserID int64, balance float64, operateType int) *Order {
	return &Order{
		UserID:   userID,
		ToUserID: toUserID,
		Balance:  balance,
		Type:     operateType,
		Status:   StatusCreate,
	}
}

func (order *Order) SaveOrderStatus(DB *gorm.DB, orderID uint, status int) error {
	err := DB.Debug().Model(order).Where("id = ?", orderID).Update("status", status).Error
	if err != nil {
		return err
	}
	return nil
}

func findByUserID(DB *gorm.DB, userID int64, limit, htype int) ([]Order, int64, error) {
	histories := make([]Order, 0, limit)
	var count *int64
	err := DB.Where("user_id = ? AND type = ?", userID, htype).
		Order("created_at DESC").
		Find(histories).Count(count).Limit(limit).Error

	return histories, *count, err
}

func findByToUserId(DB *gorm.DB, userID int64, limit, htype int) ([]Order, int64, error) {
	histories := make([]Order, 0, limit)
	var count *int64
	err := DB.Where("to_user_id = ? AND type = ?", userID, htype).
		Order("created_at DESC").
		Find(histories).Count(count).Limit(limit).Error
	return histories, *count, err
}

// GetIncrByUserID 获取用户充值历史记录
func GetIncrByUserID(DB *gorm.DB, userID int64, limit int) ([]Order, int64, error) {
	histories, count, err := findByUserID(DB, userID, limit, TypeIncr)
	if err != nil {
		return nil, 0, err
	}

	return histories, count, nil
}

// GetDecrByUserID 获取用户取款历史记录
func GetDecrByUserID(DB *gorm.DB, userID int64, limit int) ([]Order, int64, error) {
	histories, count, err := findByUserID(DB, userID, limit, TypeDecr)
	if err != nil {
		return nil, 0, err
	}

	return histories, count, nil
}

// GetTransferyOnceByUserID 获取用户单项转账历史记录
// isToMe true:转账给我的 false:我转账给别人的
func GetTransferyOnceByUserID(DB *gorm.DB, userID int64, limit int, isToMe bool) ([]Order, int64, error) {
	histories := make([]Order, 0, limit)
	var count int64
	var err error
	if isToMe {
		histories, count, err = findByUserID(DB, userID, limit, TypeTransfer)
	} else {
		histories, count, err = findByToUserId(DB, userID, limit, TypeTransfer)
	}
	if err != nil {
		return nil, 0, err
	}

	return histories, count, nil
}

// GetTransferyAllByUserID 获取所有关于我的历史记录
func GetTransferyAllByUserID(DB *gorm.DB, userID, toUserID int64, limit int) ([]Order, int64, error) {
	histories := make([]Order, 0, limit)
	var count *int64
	err := DB.Where("(user_id = ? OR to_user_id = ?) AND type = ?", userID, toUserID, TypeTransfer).
		Order("created_at DESC").
		Find(histories).Count(count).Limit(limit).Error
	if err != nil {
		return nil, 0, err
	}

	return histories, *count, nil
}
