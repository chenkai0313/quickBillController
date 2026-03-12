package response

import (
	"time"

	"github.com/shopspring/decimal"
)

type UserListResponseData struct {
	Id          int64           `json:"id"`
	CardNumber  string          `json:"card_number"`
	AliasNumber string          `json:"alias_number"`
	Status      int             `json:"status"`
	Balance     decimal.Decimal `json:"balance"` //balance
	PhoneNumber string          `json:"phone_number"`
	DestroyAt   time.Time       `json:"destroy_at"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type UserCreateResponse struct {
	Id          int64     `json:"id"`
	AliasNumber string    `json:"alias_number"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserReadCardsResponse struct {
	Id             int64                            `json:"id"` //if id is 0, it means the card is not used
	CardNumber     string                           `json:"card_number"`
	BindType       int                              `json:"bind_type"`
	PhoneNumber    string                           `json:"phone_number"`
	Status         int                              `json:"status"`          //0: destroyed, 1: using
	DestroyAt      time.Time                        `json:"destroy_at"`      //destroy time
	AliasNumber    string                           `json:"alias_number"`    //alias number
	Balance        decimal.Decimal                  `json:"balance"`         //balance
	CreatedAt      time.Time                        `json:"created_at"`      //create time
	UpdatedAt      time.Time                        `json:"updated_at"`      //update time
	BillList       []UserBillListResponseData       `json:"bill_list"`       //bill list
	TopUpList      []UserTopUpListResponseData      `json:"top_up_list"`     //top up list
	WithdrawalList []UserWithdrawalListResponseData `json:"withdrawal_list"` //withdrawal list
}
type UserWithdrawalListResponseData struct {
	Id        int64           `json:"id"`         //withdrawal id
	Amount    decimal.Decimal `json:"amount"`     //amount
	Balance   decimal.Decimal `json:"balance"`    //balance
	CreatedAt time.Time       `json:"created_at"` //create time
}

type UserBillListResponseData struct {
	Id           int64           `json:"id"`           //bill id
	OrderNumber  string          `json:"order_number"` //order number
	Amount       decimal.Decimal `json:"amount"`       //amount
	Balance      decimal.Decimal `json:"balance"`      //balance
	MerchantId   int64           `json:"merchant_id"`
	MerchantName string          `json:"merchant_name"`
	AfterBalance decimal.Decimal `json:"after_balance"` //after balance
	CreatedAt    time.Time       `json:"created_at"`    //create time
}

type UserBillSummaryResponse struct {
	Balance decimal.Decimal `json:"balance"` //balance
}

type UserTopUpListResponseData struct {
	Id        int64           `json:"id"`         //top up id
	Amount    decimal.Decimal `json:"amount"`     //amount
	Balance   decimal.Decimal `json:"balance"`    //balance
	CreatedAt time.Time       `json:"created_at"` //create time
	UpdatedAt time.Time       `json:"updated_at"` //update time
}
