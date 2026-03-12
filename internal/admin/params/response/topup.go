package response

import (
	"time"

	"github.com/shopspring/decimal"
)

type TopupRecordListResponseData struct {
	Id          int64           `json:"id"`
	UserId      int64           `json:"user_id"`
	AliasNumber string          `json:"alias_number"`
	Amount      decimal.Decimal `json:"amount"`
	Balance     decimal.Decimal `json:"balance"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
