package services

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"quickBillController/app"
	"quickBillController/config"
	"quickBillController/internal/admin/params/request"
	"quickBillController/internal/admin/params/response"
	"quickBillController/models"
)

type MerchantService struct {
}

func NewMerchantService() *MerchantService {
	return &MerchantService{}
}

func (merchant *MerchantService) MerchantBillSummary(request request.MerchantBillSummaryRequest) (resp *response.MerchantBillSummary, err error) {
	merchantModel := models.Merchant{}
	merchantBalanceData, err := merchantModel.GetMerchantBalance(request.MerchantId)
	if err != nil {
		return nil, err
	}
	return &response.MerchantBillSummary{
		FeeRate:                      merchantBalanceData.FeeRate,
		FrozenAmount:                 merchantBalanceData.FrozenAmount,
		OriginalFrozenAmount:         merchantBalanceData.OriginalFrozenAmount,
		TotalAmount:                  merchantBalanceData.TotalAmount,
		CanWithdrawalBalance:         merchantBalanceData.CanWithdrawalBalance,
		OriginalCanWithdrawalBalance: merchantBalanceData.OriginalCanWithdrawalBalance,
	}, nil
}

func (merchant *MerchantService) MerchantUpdate(request request.MerchantUpdateRequest) (resp *response.MerchantUpdateResponse, err error) {
	merchantModelOld := models.Merchant{}
	if err := merchantModelOld.GetByUserName(request.UserName); err == nil && merchantModelOld.Id != 0 && merchantModelOld.Id != request.Id {
		return nil, errors.New("merchant user name already exists")
	}
	merchantModel := models.Merchant{}
	if err := merchantModel.GetByID(request.Id); err != nil {
		return nil, err
	}
	if request.Password != "" {
		pwd, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		merchantModel.Password = string(pwd)
	}

	merchantModel.FeeRate = request.FeeRate
	merchantModel.UserName = request.UserName
	merchantModel.UpdatedAt = time.Now()
	if err := merchantModel.Save(); err != nil {
		return nil, err
	}
	return &response.MerchantUpdateResponse{
		Id:        merchantModel.Id,
		FeeRate:   merchantModel.FeeRate,
		UserName:  merchantModel.UserName,
		CreatedAt: merchantModel.CreatedAt,
		UpdatedAt: merchantModel.UpdatedAt,
	}, nil
}

func (merchant *MerchantService) MerchantList(request request.MerchantListRequest) (resp []response.MerchantListResponseData, cnt int64, err error) {
	merchantModel := models.Merchant{}
	db := app.DB.Model(&merchantModel).Order("id DESC")
	if request.UserName != "" {
		db = db.Where("user_name iLIKE ?", "%"+request.UserName+"%")
	}
	var list []models.Merchant

	if err := db.Count(&cnt).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Scopes(models.ScopePaginate(request.Page, request.PageSize)).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	for _, ll := range list {
		mm := models.Merchant{}
		merchantBalanceData, err := mm.GetMerchantBalance(ll.Id)
		if err != nil {
			return nil, 0, err
		}
		withdrawalModel := models.Withdrawal{}
		totalWithdrawalAmount, err := withdrawalModel.GetTotalWithdrawalAmountByMerchantID(ll.Id)
		if err != nil {
			return nil, 0, err
		}
		resp = append(resp, response.MerchantListResponseData{
			Id:                           ll.Id,
			UserName:                     ll.UserName,
			CreatedAt:                    ll.CreatedAt,
			UpdatedAt:                    ll.UpdatedAt,
			FeeRate:                      ll.FeeRate,
			FrozenAmount:                 merchantBalanceData.FrozenAmount,
			OriginalFrozenAmount:         merchantBalanceData.OriginalFrozenAmount,
			CanWithdrawalBalance:         merchantBalanceData.CanWithdrawalBalance,
			OriginalCanWithdrawalBalance: merchantBalanceData.OriginalCanWithdrawalBalance,
			TotalAmount:                  merchantBalanceData.TotalAmount,
			TotalWithdrawalAmount:        totalWithdrawalAmount,
		})
	}
	return resp, cnt, nil
}

func (merchant *MerchantService) MerchantCreate(request request.MerchantCreateRequest) (resp *response.MerchantCreateResponse, err error) {
	merchantModelOld := models.Merchant{}
	if err := merchantModelOld.GetByUserName(request.UserName); err == nil && merchantModelOld.Id != 0 {
		return nil, errors.New("merchant user name already exists")
	}
	defaultPassword, err := bcrypt.GenerateFromPassword([]byte(config.GetCfg().Merchant.DefaultPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	merchantModel := models.Merchant{
		UserName: request.UserName,
		FeeRate:  request.FeeRate,
		Password: string(defaultPassword),
	}
	if err := merchantModel.Insert(); err != nil {
		return nil, err
	}
	return &response.MerchantCreateResponse{
		Id:        merchantModel.Id,
		FeeRate:   merchantModel.FeeRate,
		UserName:  merchantModel.UserName,
		CreatedAt: merchantModel.CreatedAt,
		UpdatedAt: merchantModel.UpdatedAt,
	}, nil
}
