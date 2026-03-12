package response

import (
	"time"

	"github.com/shopspring/decimal"
)

type WithdrawalRecordListResponseData struct {
	Id               int64           `json:"id"`
	UserId           int64           `json:"user_id" `
	MerchantId       int64           `json:"merchant_id" `
	Amount           decimal.Decimal `json:"amount" `
	Fee              decimal.Decimal `json:"fee" `
	Balance          decimal.Decimal `json:"balance" `           //剩余可提现金额
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}


type WithdrawalSendCodeResponse struct {
	Url string `json:"url"`
}