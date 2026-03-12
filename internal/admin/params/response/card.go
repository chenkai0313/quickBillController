package response

import (
	"time"
)

type CardListResponseData struct {
	Id          int64     `json:"id"`
	Number      string    `json:"number"`
	AliasNumber string    `json:"alias_number"`
	UserId      int64     `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}