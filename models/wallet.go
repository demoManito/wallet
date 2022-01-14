package models

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"wallet/pkg/work"
	"wallet/util"

	"gorm.io/gorm"
)

const WalletManage = "wallet/manage"

var (
	InsufficientBalanceErr = errors.New("账户余额不足")
)

type Wallet struct {
	gorm.Model

	UserID  int64   `gorm:"column:user_id;unique_index:uniq_ux_userid" json:"user_id"`
	Balance float64 `gorm:"column:balance" json:"balance"`
}

func (*Wallet) TableName() string {
	return "wallets"
}

func (w *Wallet) UpdateWallet(DB *gorm.DB, userID int64, balance float64) error {
	return DB.Model(w).Where("user_id", userID).Update("balance", balance).Error
}

// FindWalletByUserID
func FindWalletByUserID(DB *gorm.DB, userID int64) (*Wallet, error) {
	wallet := new(Wallet)
	db := DB.First(wallet, Wallet{UserID: userID})
	if err := db.Error; err != nil {
		return nil, err
	}
	return wallet, nil
}

// SaveMoney 存钱
func SaveMoney(DB *gorm.DB, userID int64, balance float64) (*Order, error) {
	order := NewCreateOrder(userID, 0, balance, TypeIncr)
	err := DB.Create(order).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}

// DrawMoney 取钱
func DrawMoney(DB *gorm.DB, userID int64, balance float64) (*Order, error) {
	wallet, err := FindWalletByUserID(DB, userID)
	if err != nil {
		return nil, err
	}
	// pretreatment filter some unqualified operations first
	if wallet.Balance-balance < 0 {
		return nil, InsufficientBalanceErr
	}

	order := NewCreateOrder(userID, 0, balance, TypeDecr)
	err = DB.Create(order).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}

// Transfer 转账
func Transfer(DB *gorm.DB, userID, toUserID int64, balance float64) (*Order, error) {
	walletDecr, err := FindWalletByUserID(DB, userID)
	if err != nil {
		return nil, err
	}
	err = walletDecr.IsBalanceEnough(balance)
	if err != nil {
		return nil, err
	}
	order := NewCreateOrder(userID, toUserID, balance, TypeTransfer)
	err = DB.Create(order).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}

// 对单条记录加锁
func ForUpdateWallet(DB *gorm.DB, id int64) (*Wallet, error) {
	wallet := new(Wallet)
	err := DB.Raw("SELECT * FROM wallets WHERE id = ? FOR UPDATE", id).First(wallet).Error
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (w *Wallet) IsBalanceEnough(balance float64) error {
	if w.Balance-balance < 0 {
		return InsufficientBalanceErr
	}
	return nil
}

func FormatWalletManage(val interface{}) string {
	query := WalletManage
	if !reflect.ValueOf(val).IsZero() {
		query = fmt.Sprintf("%s/%v", WalletManage, val)
	}

	return query
}

func FormatWalletUserID(userID int64) string {
	strUserID := strconv.Itoa(int(userID))
	if !util.Includes(len(work.LargeVolume), func(i int) bool {
		if work.LargeVolume[i] == strUserID {
			return true
		}
		return false
	}) {
		strUserID = ""
	}
	return strUserID
}
