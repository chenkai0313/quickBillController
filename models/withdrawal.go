package models

import (
	"time"

	"github.com/shopspring/decimal"
	"quickBillController/app"
)

type Withdrawal struct {
	Id         int64           `gorm:"primarykey" json:"id"`
	UserId     int64           `json:"user_id" gorm:"type:int8;not null;default:0;comment:'用户ID'" `
	MerchantId int64           `json:"merchant_id" gorm:"type:int8;not null;default:0;comment:'商户ID'" `
	Amount     decimal.Decimal `json:"amount" gorm:"type:numeric(36,18);not null;comment:'申请提现金额'" `   //申请提现金额 扣除手续费之内的金额
	Fee        decimal.Decimal `json:"fee" gorm:"type:numeric(36,18);not null;comment:'提现手续费'" `       //提现手续费
	Balance    decimal.Decimal `json:"balance" gorm:"type:numeric(36,18);not null;comment:'剩余可提现金额'" ` //剩余可提现金额
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

func (m *Withdrawal) GetListByMerchantID(merchantID int64) (list []Withdrawal, err error) {
	if err := app.DB.Where("merchant_id = ?", merchantID).Order("id DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (m *Withdrawal) Insert() error {
	return app.DB.Create(m).Error
}

func (m *Withdrawal) GetListByUserId(userId int64) (list []Withdrawal, err error) {
	if err := app.DB.Where("user_id = ?", userId).Order("id DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (m *Withdrawal) GetTotalWithdrawalAmountByMerchantID(merchantID int64) (total decimal.Decimal, err error) {
	var result struct {
		Total decimal.Decimal `gorm:"column:sum"`
	}
	if err := app.DB.Model(&Withdrawal{}).
		Where("merchant_id = ?", merchantID).
		Select("COALESCE(SUM(amount), 0) as sum").
		Scan(&result).Error; err != nil {
		return decimal.Zero, err
	}
	return result.Total, nil
}
