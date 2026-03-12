package request

type CardListRequest struct {
	Page        int    `form:"page" binding:"required"`
	PageSize    int    `form:"page_size" binding:"required"`
	Number      string `form:"number" `
	AliasNumber string `form:"alias_number" `
	UserId      int64  `form:"user_id" `
}
