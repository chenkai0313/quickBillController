package response

import (
	"time"

	"github.com/shopspring/decimal"
)

type BillRecordListResponseData struct {
	Id          int64           `json:"id"`           //bill id
	OrderNumber string          `json:"order_number"` //order number
	UserId      int64           `json:"user_id" `
	PhoneNumber string          `json:"phone_number" `
	MerchantId  int64           `json:"merchant_id" `
	MerchantName string          `json:"merchant_name" `
	AliasNumber string          `json:"alias_number" `
	Amount      decimal.Decimal `json:"amount" `
	Balance     decimal.Decimal `json:"balance" `
	CreatedAt   time.Time       `json:"created_at"`
}
