package request

import "github.com/shopspring/decimal"

type UserListRequest struct {
	Page        int    `form:"page" binding:"required"`
	PageSize    int    `form:"page_size" binding:"required"`
	AliasNumber string `form:"alias_number"`
	PhoneNumber string `form:"phone_number"`
}

type UserCreateRequest struct {
	CardNumber  string          `json:"card_number" binding:"required"`
	AliasNumber string          `json:"alias_number" binding:"required"`
	BindType    int             `json:"bind_type" ` // 0:null 1:password 2:phone
	Password    string          `json:"password"`
	PhoneNumber string          `json:"phone_number"`
	TopAmount   decimal.Decimal `json:"top_amount"`
}

type UserReadCardsRequest struct {
	CardNumber string `form:"card_number" binding:"required"`
}

type UserBillSummaryRequest struct {
	CardNumber string `form:"card_number" binding:"required"`
}

type UserBillRecordListRequest struct {
	Page       int    `form:"page" binding:"required"`
	PageSize   int    `form:"page_size" binding:"required"`
	CardNumber string `form:"card_number" binding:"required"`
}
