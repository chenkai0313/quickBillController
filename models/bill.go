package models

import (
	"fmt"
	"time"

	"quickBillController/app"
	"quickBillController/utils"

	"github.com/shopspring/decimal"
)

type Bill struct {
	Id           int64           `gorm:"primarykey" json:"id"`
	OrderNumber  string          `json:"order_number" gorm:"type:varchar(1000);not null;unique"`
	UserId       int64           `json:"user_id" gorm:"type:int8;not null;default:0;index"`
	MerchantId   int64           `json:"merchant_id" gorm:"type:int8;not null;default:0;index"`
	CardNumber   string          `json:"card_number" gorm:"type:varchar(1000);not null;index;default:''" `
	AliasNumber  string          `gorm:"type:varchar(1000);not null;index;default:''"  json:"alias_number"`
	Amount       decimal.Decimal `json:"amount" gorm:"type:numeric(36,18);not null;default:0;comment:'消费金额'" `         //consumption amount
	Balance      decimal.Decimal `json:"balance" gorm:"type:numeric(36,18);not null;default:0;comment:'当前余额'" `        //current balance
	AfterBalance decimal.Decimal `json:"after_balance" gorm:"type:numeric(36,18);not null;default:0;comment:'消费后余额'" ` //after balance
	CreatedAt    time.Time       `json:"created_at"`
	MerchantName string          `json:"merchant_name" gorm:"-"` //merchant name
}

func (m *Bill) Insert() error {
	return app.DB.Create(m).Error
}

func (m *Bill) GetByUserId(userId int64) error {
	if err := app.DB.Where("user_id = ?", userId).Order("id DESC").First(m).Error; err != nil {
		return err
	}
	return nil
}

// 联表查询 查询出商户名称
func (m *Bill) GetListByUserId(userId int64) (list []Bill, err error) {
	type BillWithMerchant struct {
		Bill
		MerchantName string `gorm:"column:merchant_name"`
	}

	var results []BillWithMerchant
	if err := app.DB.Table("sbs_bills").
		Select("sbs_bills.*, sbs_merchants.user_name as merchant_name").
		Joins("LEFT JOIN sbs_merchants ON sbs_bills.merchant_id = sbs_merchants.id").
		Where("sbs_bills.user_id = ?", userId).
		Order("sbs_bills.id DESC").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	// 将结果转换为 Bill 列表
	list = make([]Bill, len(results))
	for i, r := range results {
		list[i] = r.Bill
		list[i].MerchantName = r.MerchantName
	}

	return list, nil
}

// 时间戳+userid+merchantid+3位随机数
// userid 为 4 位数字 不足 4 位前面补 0
// merchantid 为 3 位数字 不足 3 位前面补 0

func (m *Bill) GenerateOrderNumber(userID int64, merchantID int64) string {
	userIdStr := fmt.Sprintf("%04d", userID)
	merchantIdStr := fmt.Sprintf("%03d", merchantID)
	randomStr := fmt.Sprintf("%03d", utils.RandInt(1000))
	return fmt.Sprintf("%d%s%s%s", time.Now().Unix(), userIdStr, merchantIdStr, randomStr)
}

func (m *Bill) GetListByMerchantID(merchantID int64) (list []Bill, err error) {
	if err := app.DB.Where("merchant_id = ?", merchantID).Order("id DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (m *Bill) GetByOrderNumber(orderNumber string) error {
	if err := app.DB.Where("order_number = ?", orderNumber).First(m).Error; err != nil {
		return err
	}
	return nil
}
