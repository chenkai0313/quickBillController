package request

import "github.com/shopspring/decimal"

type MerchantLoginRequest struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type MerchantUndoBillRecordListRequest struct {
	Page        int    `form:"page" binding:"required"`
	PageSize    int    `form:"page_size" binding:"required"`
	OrderNumber string `form:"order_number" `
	StartAt     int64  `form:"start_at" `
	EndAt       int64  `form:"end_at" `
}

type MerchantBillRecordListRequest struct {
	Page        int    `form:"page" binding:"required"`
	PageSize    int    `form:"page_size" binding:"required"`
	OrderNumber string `form:"order_number" `
	StartAt     int64  `form:"start_at" `
	EndAt       int64  `form:"end_at" `
}

type MerchantConsumeRequest struct {
	CardNumber string          `json:"card_number" binding:"required"`
	Amount     decimal.Decimal `json:"amount" binding:"required"`
}

type MerchantUndoConsumeRequest struct {
	OrderNumber string `json:"order_number" binding:"required"`
}

type MerchantUpdatePasswordRequest struct {
	OriginalPassword string `json:"original_password" binding:"required"`
	Password         string `json:"password" binding:"required"`
}

type MerchantQueryCardBalanceRequest struct {
	CardNumber string `form:"card_number" binding:"required"`
}
