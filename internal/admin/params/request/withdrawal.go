package request

import "github.com/shopspring/decimal"

type WithdrawalRequest struct {
	UserId     int64           `json:"user_id"`
	MerchantId int64           `json:"merchant_id" `
	Amount     decimal.Decimal `json:"amount" binding:"required"`
	Password   string          `json:"password"`
	Code       string          `json:"code"`
}

type WithdrawalSendCodeRequest struct {
	UserId int64 `json:"user_id" binding:"required"`
}
