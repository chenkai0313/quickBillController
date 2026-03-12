package response

import (
	"time"

	"github.com/shopspring/decimal"
)

type MerchantLoginResponse struct {
	Id        int64  `json:"id"`
	UserName  string `json:"user_name"`
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

type MerchantBillRecordListResponseData struct {
	Id           int64           `json:"id"`           //bill id
	OrderNumber  string          `json:"order_number"` //order number
	UserId       int64           `json:"user_id" `
	MerchantId   int64           `json:"merchant_id" `
	AliasNumber  string          `json:"alias_number" `
	Amount       decimal.Decimal `json:"amount" `
	AfterBalance decimal.Decimal `json:"after_balance" `
	Balance      decimal.Decimal `json:"balance" `
	CreatedAt    time.Time       `json:"created_at"`
}

type MerchantSummaryResponse struct {
	FeeRate              decimal.Decimal `json:"fee_rate"`              // fee rate
	FrozenAmount         decimal.Decimal `json:"frozen_amount"`          //total balance
	OriginalFrozenAmount decimal.Decimal `json:"original_frozen_amount"` //original total balance
	CanWithdrawalBalance decimal.Decimal `json:"can_withdrawal_balance"` //can withdrawal balance
	OriginalCanWithdrawalBalance decimal.Decimal `json:"original_can_withdrawal_balance"` //original can withdrawal balance
	TotalAmount                  decimal.Decimal `json:"total_amount"`                    // total amount
	TotalWithdrawalAmount         decimal.Decimal `json:"total_withdrawal_amount"`          // 总提现金额
}

type MerchantConsumeResponse struct {
	Id           int64           `json:"id"`
	OrderNumber  string          `json:"order_number"`
	UserId       int64           `json:"user_id"`
	MerchantId   int64           `json:"merchant_id"`
	AliasNumber  string          `json:"alias_number"`
	Amount       decimal.Decimal `json:"amount"`
	Balance      decimal.Decimal `json:"balance"`
	AfterBalance decimal.Decimal `json:"after_balance"`
	CreatedAt    time.Time       `json:"created_at"`
}

type MerchantQueryCardBalanceResponse struct {
	Balance decimal.Decimal `json:"balance"`
}
