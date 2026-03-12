package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type UndoBill struct {
	Id           int64           `gorm:"primarykey" json:"id"`
	OrderNumber  string          `json:"order_number" gorm:"type:varchar(1000);not null;unique"`
	UserId       int64           `json:"user_id" gorm:"type:int8;not null;default:0;index"`
	MerchantId   int64           `json:"merchant_id" gorm:"type:int8;not null;default:0;index"`
	CardNumber   string          `json:"card_number" gorm:"type:varchar(1000);not null;index;default:''" `
	AliasNumber  string          `gorm:"type:varchar(1000);not null;index;default:''"  json:"alias_number"`
	Amount       decimal.Decimal `json:"amount" gorm:"type:numeric(36,18);not null;default:0;comment:'消费金额'" `         //consumption amount
	AfterBalance decimal.Decimal `json:"after_balance" gorm:"type:numeric(36,18);not null;default:0;comment:'消费后余额'" ` //after balance
	Balance      decimal.Decimal `json:"balance" gorm:"type:numeric(36,18);not null;default:0;comment:'当前余额'" `        //current balance
	CreatedAt    time.Time       `json:"created_at"`
}
