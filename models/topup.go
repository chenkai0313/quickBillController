package models

import (
	"time"

	"github.com/shopspring/decimal"
	"quickBillController/app"
)

type Topup struct {
	Id          int64           `gorm:"primarykey" json:"id"`
	UserId      int64           `gorm:"type:int8;not null;default:0;comment:'用户ID'" json:"user_id"`
	CardNumber  string          `gorm:"type:varchar(1000);not null;comment:'充值卡号'" json:"card_number"` //充值卡号
	AliasNumber string          `gorm:"type:varchar(1000);not null;index;default:''"  json:"alias_number"`
	Amount      decimal.Decimal `gorm:"type:numeric(36,18);not null;default:0;comment:'充值金额'" json:"amount"`   //充值金额
	Balance     decimal.Decimal `gorm:"type:numeric(36,18);not null;default:0;comment:'充值签余额'" json:"balance"` //充值前余额
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func (m *Topup) Insert() error {
	return app.DB.Create(m).Error
}

func (m *Topup) GetListByUserId(userId int64) (list []Topup, err error) {
	if err := app.DB.Where("user_id = ?", userId).Order("id DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
