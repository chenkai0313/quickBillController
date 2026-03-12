package request

import "github.com/shopspring/decimal"

type MerchantListRequest struct {
	Page     int    `form:"page" binding:"required"`
	PageSize int    `form:"page_size" binding:"required"`
	UserName string `form:"user_name" `
}

type MerchantUpdateRequest struct {
	Id       int64           `json:"id" binding:"required"`
	UserName string          `json:"user_name" binding:"required"`
	Password string          `json:"password" binding:"required"`
	FeeRate  decimal.Decimal `json:"fee_rate" binding:"required"`
}

type MerchantCreateRequest struct {
	UserName string          `json:"user_name" binding:"required"`
	FeeRate  decimal.Decimal `json:"fee_rate" binding:"required"`
}

type MerchantBillSummaryRequest struct {
	MerchantId int64 `form:"merchant_id" binding:"required"`
}

type WithdrawalRecordListRequest struct {
	Page       int    `form:"page" binding:"required"`
	PageSize   int    `form:"page_size" binding:"required"`
	MerchantId int64  `form:"merchant_id" `
	UserID     int64  `form:"user_id" `
	CardNumber string `form:"card_number" `
}

type WithdrawalRecordListExportRequest struct {
	MerchantId int64  `form:"merchant_id" `
	UserID     int64  `form:"user_id" `
	CardNumber string `form:"card_number" `
}
