package request

import "github.com/shopspring/decimal"

type TopupRequest struct {
	CardNumber string          `json:"card_number"` //充值卡号
	Amount     decimal.Decimal `json:"amount"`      //充值金额
}

type TopupRecordListRequest struct {
	Page        int    `form:"page" binding:"required"`
	PageSize    int    `form:"page_size" binding:"required"`
	AliasNumber string `form:"alias_number"`
	UserID      int64  `form:"user_id"`
}

type TopupRecordListExportRequest struct {
	AliasNumber string `form:"alias_number"`
	UserID      int64  `form:"user_id"`
}
