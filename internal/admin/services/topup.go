package services

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"quickBillController/app"
	"quickBillController/internal/admin/params/request"
	"quickBillController/internal/admin/params/response"
	"quickBillController/models"

	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"
	"quickBillController/utils/carrier"
)

type TopupService struct {
}

func NewTopupService() *TopupService {
	return &TopupService{}
}

func (topup *TopupService) TopupRecordListExport(request request.TopupRecordListExportRequest) (f *excelize.File, fn string, err error) {
	topupModel := models.Topup{}
	db := app.DB.Model(&topupModel).Order("id DESC")

	if request.AliasNumber != "" {
		db = db.Where("alias_number = ?", request.AliasNumber)
	}
	if request.UserID != 0 {
		db = db.Where("user_id = ?", request.UserID)
	}
	list := []models.Topup{}
	if err := db.Find(&list).Error; err != nil {
		return nil, "", err
	}
	var car carrier.ExcelCarrier
	car.Titles = []string{
		"id",
		"Card Number",
		"Alias Number",
		"Amount",
		"Balance",
		"Created At",
	}
	car.File = excelize.NewFile()
	car.SheetName = "Sheet1"
	car.Data = make([][]string, len(list))
	for i, t := range list {
		car.Data[i] = append(car.Data[i], cast.ToString(t.Id))
		car.Data[i] = append(car.Data[i], t.CardNumber)
		car.Data[i] = append(car.Data[i], t.AliasNumber)
		car.Data[i] = append(car.Data[i], t.Amount.String())
		car.Data[i] = append(car.Data[i], t.Balance.String())
		car.Data[i] = append(car.Data[i], t.CreatedAt.Format("2006-01-02 15:04:05"))
	}
	// Write data to Excel file
	if err := car.Write(); err != nil {
		return nil, "", err
	}
	// Return file with proper filename
	filename := "topup_records.xlsx"
	return car.File, filename, nil
}

func (topup *TopupService) TopupRecordList(request request.TopupRecordListRequest) (resp []response.TopupRecordListResponseData, cnt int64, err error) {
	topupModel := models.Topup{}

	db := app.DB.Model(&topupModel).Order("id DESC")

	if request.AliasNumber != "" {
		cardModel := models.Card{}
		if err := cardModel.GetByAliasNumber(request.AliasNumber); err != nil {
			return resp, 0, nil
		}
		userModel := models.User{}
		if err := userModel.GetById(cardModel.UserId); err != nil {
			return resp, 0, nil
		}
		if userModel.Status == models.UserStatusDestroy {
			return resp, 0, nil
		}
		db = db.Where("user_id = ?", userModel.Id)
	}

	if request.UserID != 0 {
		db = db.Where("user_id = ?", request.UserID)
	}

	if err := db.Count(&cnt).Error; err != nil {
		return nil, 0, err
	}
	list := []models.Topup{}
	if err := db.Scopes(models.ScopePaginate(request.Page, request.PageSize)).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	for _, ll := range list {
		resp = append(resp, response.TopupRecordListResponseData{
			Id:          ll.Id,
			UserId:      ll.UserId,
			AliasNumber: ll.AliasNumber,
			Amount:      ll.Amount,
			Balance:     ll.Balance,
			CreatedAt:   ll.CreatedAt,
			UpdatedAt:   ll.UpdatedAt,
		})
	}
	return resp, cnt, nil
}

func (topup *TopupService) Topup(request request.TopupRequest) (err error) {
	if request.Amount.IsZero() {
		return fmt.Errorf("invalid amount")
	}
	cardModel := models.Card{}
	if err := cardModel.GetByCardNumber(request.CardNumber); err != nil {
		app.ZapLog.Error("card number not found", zap.Error(err))
		return fmt.Errorf("card number not found")
	}

	userModel := models.User{}
	if err := userModel.GetById(cardModel.UserId); err != nil {
		app.ZapLog.Error("user not found", zap.Error(err))
		return fmt.Errorf("user not found")
	}
	if userModel.Status == models.UserStatusDestroy {
		app.ZapLog.Error("user is destroyed", zap.Error(err))
		return fmt.Errorf("user is destroyed")
	}
	userBalance, err := userModel.GetUserBalance(cardModel.UserId)
	if err != nil {
		app.ZapLog.Error("user balance not found", zap.Error(err))
		return fmt.Errorf("user balance not found")
	}

	topupModel := models.Topup{
		UserId:      userModel.Id,
		CardNumber:  request.CardNumber,
		AliasNumber: cardModel.AliasNumber,
		Amount:      request.Amount,
		Balance:     userBalance,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := topupModel.Insert(); err != nil {
		app.ZapLog.Error("create topup failed", zap.Error(err))
		return fmt.Errorf("create topup failed %v", err)
	}
	return nil
}
