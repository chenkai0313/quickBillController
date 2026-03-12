package response

import (
	"time"

	"github.com/shopspring/decimal"
)

type MerchantListResponseData struct {
	Id                           int64           `json:"id"`
	UserName                     string          `json:"user_name"`
	CreatedAt                    time.Time       `json:"created_at"`
	UpdatedAt                    time.Time       `json:"updated_at"`
	FeeRate                      decimal.Decimal `json:"fee_rate"`                        // 手续费率
	FrozenAmount                 decimal.Decimal `json:"frozen_amount"`                   //扣除手续费之后的冻结金额
	OriginalFrozenAmount         decimal.Decimal `json:"original_frozen_amount"`          // 原始冻结金额 (24小时内的订单数据 不算入)
	TotalAmount                  decimal.Decimal `json:"total_amount"`                    // 总营业额包括冻结时间内的订单
	CanWithdrawalBalance         decimal.Decimal `json:"can_withdrawal_balance"`          //扣除手续费之后的可提现金额
	OriginalCanWithdrawalBalance decimal.Decimal `json:"original_can_withdrawal_balance"` // 原始可提现金额
	TotalWithdrawalAmount         decimal.Decimal `json:"total_withdrawal_amount"`          // 总提现金额
}

type MerchantCreateResponse struct {
	Id        int64           `json:"id"`
	UserName  string          `json:"user_name"`
	FeeRate   decimal.Decimal `json:"fee_rate"` // 手续费率
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type MerchantUpdateResponse struct {
	Id        int64           `json:"id"`
	UserName  string          `json:"user_name"`
	FeeRate   decimal.Decimal `json:"fee_rate"` // 手续费率
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type MerchantBillSummary struct {	
	FeeRate                      decimal.Decimal `json:"fee_rate"`                        // 手续费率
	FrozenAmount                 decimal.Decimal `json:"frozen_amount"`                   //扣除手续费之后的冻结金额
	OriginalFrozenAmount         decimal.Decimal `json:"original_frozen_amount"`          // 原始冻结金额 (24小时内的订单数据 不算入)
	TotalAmount                  decimal.Decimal `json:"total_amount"`                    // 总营业额包括冻结时间内的订单
	CanWithdrawalBalance         decimal.Decimal `json:"can_withdrawal_balance"`          //扣除手续费之后的可提现金额
	OriginalCanWithdrawalBalance decimal.Decimal `json:"original_can_withdrawal_balance"` // 原始可提现金额
}
