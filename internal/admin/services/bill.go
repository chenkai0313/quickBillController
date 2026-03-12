package services

import (
	"quickBillController/app"
	"quickBillController/internal/admin/params/request"
	"quickBillController/internal/admin/params/response"
	"quickBillController/models"
	"quickBillController/utils/carrier"

	"github.com/xuri/excelize/v2"

	"github.com/spf13/cast"
	"time"
)

type BillService struct {
}

func NewBillService() *BillService {
	return &BillService{}
}

func (bill *BillService) BillRecordListExport(request request.BillRecordListExportRequest) (f *excelize.File, fn string, err error) {
	type BillWithMerchant struct {
		models.Bill
		PhoneNumber  string `gorm:"column:phone_number"`
		MerchantName string `gorm:"column:merchant_name"`
	}

	// 构建基础查询，使用联表查询获取商户名称
	db := app.DB.Table("sbs_bills").
		Select("sbs_bills.*, sbs_merchants.user_name as merchant_name, sbs_users.phone_number as phone_number").
		Joins("LEFT JOIN sbs_merchants ON sbs_bills.merchant_id = sbs_merchants.id").
		Joins("LEFT JOIN sbs_users ON sbs_bills.user_id = sbs_users.id").
		Order("sbs_bills.id DESC")

	// 添加查询条件
	if request.MerchantName != "" {
		db = db.Where("sbs_merchants.user_name iLIKE ?", "%"+request.MerchantName+"%")
	}
	if request.UserId != 0 {
		db = db.Where("sbs_bills.user_id = ?", request.UserId)
	}
	if request.OrderNumber != "" {
		db = db.Where("sbs_bills.order_number iLIKE ?", "%"+request.OrderNumber+"%")
	}
	if request.PhoneNumber != "" {
		db = db.Where("sbs_users.phone_number iLIKE ?", "%"+request.PhoneNumber+"%")
	}
	if request.StartTime != 0 {
		db = db.Where("sbs_bills.created_at >= ?", time.Unix(request.StartTime, 0))
	}
	if request.EndTime != 0 {
		db = db.Where("sbs_bills.created_at <= ?", time.Unix(request.EndTime, 0))
	}
	// 查询数据
	var results []BillWithMerchant
	if err := db.Scan(&results).Error; err != nil {
		return nil, "", err
	}

	var car carrier.ExcelCarrier
	car.Titles = []string{
		"id",
		"Order Number",
		"User ID",
		"Phone Number",
		"Merchant Name",
		"Alias Number",
		"Amount",
		"Balance After",
		"Created At",
	}
	car.File = excelize.NewFile()
	car.SheetName = "Sheet1"
	car.Data = make([][]string, len(results))
	for i, b := range results {
		car.Data[i] = append(car.Data[i], cast.ToString(b.Id))
		car.Data[i] = append(car.Data[i], b.OrderNumber)
		car.Data[i] = append(car.Data[i], cast.ToString(b.UserId))
		car.Data[i] = append(car.Data[i], b.PhoneNumber)
		car.Data[i] = append(car.Data[i], b.MerchantName)
		car.Data[i] = append(car.Data[i], b.AliasNumber)
		car.Data[i] = append(car.Data[i], b.Amount.String())
		car.Data[i] = append(car.Data[i], b.Balance.String())
		car.Data[i] = append(car.Data[i], b.CreatedAt.Format("2006-01-02 15:04:05"))
	}
	// Write data to Excel file
	if err := car.Write(); err != nil {
		return nil, "", err
	}
	// Return file with proper filename
	filename := "bill_records.xlsx"
	return car.File, filename, nil
}

// 联表查询 查询出商户名称
func (bill *BillService) BillRecordList(request request.BillRecordListRequest) (resp []response.BillRecordListResponseData, cnt int64, err error) {
	type BillWithMerchant struct {
		models.Bill
		PhoneNumber  string `gorm:"column:phone_number"`
		MerchantName string `gorm:"column:merchant_name"`
	}

	// 构建基础查询
	db := app.DB.Table("sbs_bills").
		Select("sbs_bills.*, sbs_merchants.user_name as merchant_name, sbs_users.phone_number as phone_number").
		Joins("LEFT JOIN sbs_users ON sbs_bills.user_id = sbs_users.id").
		Joins("LEFT JOIN sbs_merchants ON sbs_bills.merchant_id = sbs_merchants.id").
		Order("sbs_bills.id DESC")

	// 添加查询条件
	if request.MerchantName != "" {
		db = db.Where("sbs_merchants.user_name iLIKE ?", "%"+request.MerchantName+"%")
	}
	if request.UserId != 0 {
		db = db.Where("sbs_bills.user_id = ?", request.UserId)
	}
	if request.OrderNumber != "" {
		db = db.Where("sbs_bills.order_number iLIKE ?", "%"+request.OrderNumber+"%")
	}
	if request.PhoneNumber != "" {
		db = db.Where("sbs_users.phone_number iLIKE ?", "%"+request.PhoneNumber+"%")
	}
	if request.StartTime != 0 {
		db = db.Where("sbs_bills.created_at >= ?", time.Unix(request.StartTime, 0))
	}
	if request.EndTime != 0 {
		db = db.Where("sbs_bills.created_at <= ?", time.Unix(request.EndTime, 0))
	}
	// 计算总数
	if err := db.Count(&cnt).Error; err != nil {
		return nil, 0, err
	}

	// 查询数据
	var results []BillWithMerchant
	if err := db.Scopes(models.ScopePaginate(request.Page, request.PageSize)).Scan(&results).Error; err != nil {
		return nil, 0, err
	}

	// 转换为响应数据
	resp = make([]response.BillRecordListResponseData, len(results))
	for i, result := range results {
		resp[i] = response.BillRecordListResponseData{
			Id:           result.Id,
			OrderNumber:  result.OrderNumber,
			PhoneNumber:  result.PhoneNumber,
			UserId:       result.UserId,
			MerchantId:   result.MerchantId,
			MerchantName: result.MerchantName,
			AliasNumber:  result.AliasNumber,
			Amount:       result.Amount,
			Balance:      result.Balance,
			CreatedAt:    result.CreatedAt,
		}
	}

	return resp, cnt, nil
}
