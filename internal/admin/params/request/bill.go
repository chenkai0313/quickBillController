package request

type BillRecordListRequest struct {
	Page         int    `form:"page" binding:"required"`
	PageSize     int    `form:"page_size" binding:"required"`
	MerchantName string `form:"merchant_name"`
	UserId       int64  `form:"user_id"`
	OrderNumber  string `form:"order_number"`
	PhoneNumber  string `form:"phone_number"`
	StartTime    int64  `form:"start_time"`
	EndTime      int64  `form:"end_time"`
}

type BillRecordListExportRequest struct {
	MerchantName string `form:"merchant_name"`
	UserId       int64  `form:"user_id"`
	OrderNumber  string `form:"order_number"`
	PhoneNumber  string `form:"phone_number"`
	StartTime    int64  `form:"start_time"`
	EndTime      int64  `form:"end_time"`
}
