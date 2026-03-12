package services

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
	"quickBillController/app"
	"quickBillController/internal/admin/params/request"
	"quickBillController/internal/admin/params/response"
	"quickBillController/models"
)

type CardService struct {
}

func NewCardService() *CardService {
	return &CardService{}
}

func (card *CardService) CardList(request request.CardListRequest) (resp []response.CardListResponseData, cnt int64, err error) {
	cardModel := models.Card{}
	db := app.DB.Model(&cardModel).Order("id DESC")
	if request.Number != "" {
		db = db.Where("number = ?", request.Number)
	}
	if request.AliasNumber != "" {
		db = db.Where("alias_number = ?", request.AliasNumber)
	}
	if request.UserId != 0 {
		db = db.Where("user_id = ?", request.UserId)
	}
	list := []models.Card{}
	if err := db.Count(&cnt).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Scopes(models.ScopePaginate(request.Page, request.PageSize)).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	for _, card := range list {
		resp = append(resp, response.CardListResponseData{
			Id:          card.Id,
			Number:      card.Number,
			AliasNumber: card.AliasNumber,
			UserId:      card.UserId,
			CreatedAt:   card.CreatedAt,
			UpdatedAt:   card.UpdatedAt,
		})
	}
	return resp, cnt, nil
}

func (card *CardService) CardImport(c *gin.Context, f *excelize.File) (err error) {
	sn := f.GetSheetName(0)
	rows, err := f.GetRows(sn)
	if err != nil {
		app.ZapLog.Error("get data from excel sheet1 error", zap.Error(err))
		return err
	}
	if len(rows) < 1 {
		return fmt.Errorf("get empty data from second row")
	}
	for _, row := range rows[1:] {
		carModel := models.Card{}
		for i, content := range row {
			switch i {
			case 0:
				carModel.Number = content
			case 1:
				carModel.AliasNumber = content
			}
		}
		carModel.CreatedAt = time.Now()
		carModel.UpdatedAt = time.Now()
		_ = carModel.GetByCardNumber(carModel.Number)
		if carModel.Id != 0 {
			continue
		}
		if err := carModel.Insert(); err != nil {
			return err
		}
	}
	return nil
}
