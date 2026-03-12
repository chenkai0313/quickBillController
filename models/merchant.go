package models

import (
	"time"

	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
	"quickBillController/app"
	"quickBillController/config"
)

type Merchant struct {
	Id        int64           `gorm:"primarykey" json:"id"`
	UserName  string          `gorm:"type:varchar(100);not null;unique;" json:"user_name"`
	Password  string          `gorm:"type:varchar(255);not null;comment:'密码'" json:"password"`
	FeeRate   decimal.Decimal `gorm:"type:numeric(36,18);not null;default:0;comment:'手续费率'" json:"fee_rate"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

func (m *Merchant) GetByUserName(userName string) error {
	if err := app.DB.Where("user_name = ?", userName).First(m).Error; err != nil {
		return err
	}
	return nil
}

func (m *Merchant) ComparePassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(password)); err != nil {
		return err
	}
	return nil
}

func (m *Merchant) GetByID(id int64) error {
	if err := app.DB.Where("id = ?", id).First(m).Error; err != nil {
		return err
	}
	return nil
}

func (m *Merchant) Insert() error {
	return app.DB.Create(m).Error
}

func (m *Merchant) Save() error {
	return app.DB.Save(m).Error
}

// 只显示扣除手续费之后的额度
// canWithdrawalBalance 可提现金额
// FrozenAmount 冻结金额 (冻结时间内的订单数据 不算入)

type MerchantBalanceData struct {
	FrozenAmount                 decimal.Decimal //扣除手续费之后的冻结金额
	OriginalFrozenAmount         decimal.Decimal //原始冻结金额
	CanWithdrawalBalance         decimal.Decimal //扣除手续费之后的可提现金额
	OriginalCanWithdrawalBalance decimal.Decimal //原始可提现金额
	FeeRate                      decimal.Decimal //手续费率
	TotalAmount                  decimal.Decimal //商户到目前所有订单的总营业额包括冻结时间内的订单
}

func (m *Merchant) GetMerchantBalance(merchantId int64) (data MerchantBalanceData, err error) {
	if err := m.GetByID(merchantId); err != nil {
		return data, err
	}

	// 计算冻结时间前的时间点
	now := time.Now()
	frozenHours := config.GetCfg().Merchant.FrozenHours
	frozenHoursAgo := now.Add(-time.Duration(frozenHours) * time.Hour)

	// 1. 获取冻结时间之前的所有的商铺订单金额总和
	var billAmountBeforeFrozenHours struct {
		Total decimal.Decimal
	}
	if err := app.DB.Model(&Bill{}).
		Where("merchant_id = ? AND created_at < ?", merchantId, frozenHoursAgo).
		Select("COALESCE(SUM(amount), 0) as total").
		Scan(&billAmountBeforeFrozenHours).Error; err != nil {
		return data, err
	}

	// 2. 获取已经提现的金额总和（使用withdrawal_amount，即实际扣除的金额）
	var withdrawalAmount struct {
		Total    decimal.Decimal `json:"total"`
		TotalFee decimal.Decimal `json:"total_fee"`
	}
	if err := app.DB.Model(&Withdrawal{}).
		Where("merchant_id = ?", merchantId).
		Select("COALESCE(SUM(amount), 0) as total,COALESCE(SUM(fee), 0) as total_fee").
		Scan(&withdrawalAmount).Error; err != nil {
		return data, err
	}

	// 3. 获取冻结时间内的商铺订单金额总和
	var billAmountWithinFrozenHours struct {
		Total decimal.Decimal
	}
	if err := app.DB.Model(&Bill{}).
		Where("merchant_id = ? AND created_at >= ?", merchantId, frozenHoursAgo).
		Select("COALESCE(SUM(amount), 0) as total").
		Scan(&billAmountWithinFrozenHours).Error; err != nil {
		return data, err
	}

	// 4. 计算手续费率
	withdrawalFee := m.FeeRate
	feeMultiplier := decimal.NewFromInt(1).Sub(withdrawalFee) // (1 - 手续费率)

	// 5. 计算可提现金额：canWithdrawalBalance = (冻结时间前的订单金额总和 - (已提现金额总和+已经提现的手续费)) * (1 - 手续费率)
	totalBillAmountBeforeFrozenHours := billAmountBeforeFrozenHours.Total
	totalWithdrawalAmount := withdrawalAmount.Total.Add(withdrawalAmount.TotalFee)
	availableAmount := totalBillAmountBeforeFrozenHours.Sub(totalWithdrawalAmount)
	data.CanWithdrawalBalance = availableAmount.Mul(feeMultiplier)

	// 6. 计算冻结金额：frozenAmount = 冻结时间内的订单金额总和 * (1 - 手续费率)
	data.FrozenAmount = billAmountWithinFrozenHours.Total.Mul(feeMultiplier)

	data.OriginalCanWithdrawalBalance = availableAmount
	data.OriginalFrozenAmount = billAmountWithinFrozenHours.Total
	data.FeeRate = withdrawalFee
	data.TotalAmount = totalBillAmountBeforeFrozenHours.Add(billAmountWithinFrozenHours.Total)
	return data, nil
}
