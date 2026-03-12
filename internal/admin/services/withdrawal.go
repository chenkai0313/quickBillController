package services

import (
	"errors"
	"fmt"
	"time"

	"quickBillController/app"
	"quickBillController/internal/admin/params/request"
	"quickBillController/internal/admin/params/response"
	"quickBillController/models"
	"quickBillController/utils"
	"quickBillController/utils/carrier"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"
)

type WithdrawalService struct {
}

func NewWithdrawalService() *WithdrawalService {
	return &WithdrawalService{}
}

func (withdrawal *WithdrawalService) WithdrawalRecordListExport(request request.WithdrawalRecordListExportRequest) (f *excelize.File, fn string, err error) {
	withdrawalModel := models.Withdrawal{}
	db := app.DB.Model(&withdrawalModel).Order("id DESC")
	if request.MerchantId != 0 {
		db = db.Where("merchant_id = ?", request.MerchantId)
	}
	if request.UserID != 0 {
		db = db.Where("user_id = ?", request.UserID)
	}
	if request.CardNumber != "" {
		card := models.Card{}
		if err := card.GetByCardNumber(request.CardNumber); err != nil {
			return nil, "", fmt.Errorf("invalid card number: %s", err.Error())
		}
		request.UserID = card.UserId
	}
	list := []models.Withdrawal{}
	if err := db.Find(&list).Error; err != nil {
		return nil, "", err
	}
	var car carrier.ExcelCarrier
	car.Titles = []string{
		"id",
		"type",
		"User ID",
		"Merchant ID",
		"Amount",
		"Fee",
		"Balance",
		"Created At",
		"Updated At",
	}
	car.File = excelize.NewFile()
	car.SheetName = "Sheet1"
	car.Data = make([][]string, len(list))
	for i, w := range list {
		car.Data[i] = append(car.Data[i], cast.ToString(w.Id))
		if w.MerchantId > 0 {
			car.Data[i] = append(car.Data[i], "merchant")
		} else {
			car.Data[i] = append(car.Data[i], "user")
		}
		car.Data[i] = append(car.Data[i], cast.ToString(w.UserId))
		car.Data[i] = append(car.Data[i], cast.ToString(w.MerchantId))
		car.Data[i] = append(car.Data[i], w.Amount.String())
		car.Data[i] = append(car.Data[i], w.Fee.String())
		car.Data[i] = append(car.Data[i], w.Balance.String())
		car.Data[i] = append(car.Data[i], w.CreatedAt.Format("2006-01-02 15:04:05"))
		car.Data[i] = append(car.Data[i], w.UpdatedAt.Format("2006-01-02 15:04:05"))
	}
	// Write data to Excel file
	if err := car.Write(); err != nil {
		return nil, "", fmt.Errorf("write excel file failed: %s", err.Error())
	}
	// Return file with proper filename
	filename := "withdrawal_records.xlsx"
	return car.File, filename, nil
}
func (withdrawal *WithdrawalService) WithdrawalRecordList(request request.WithdrawalRecordListRequest) (resp []response.WithdrawalRecordListResponseData, cnt int64, err error) {
	withdrawalModel := models.Withdrawal{}
	db := app.DB.Model(&withdrawalModel).Order("id DESC")
	if request.MerchantId != 0 {
		db = db.Where("merchant_id = ?", request.MerchantId)
	}
	if request.CardNumber != "" {
		card := models.Card{}
		if err := card.GetByCardNumber(request.CardNumber); err != nil {
			return nil, 0, fmt.Errorf("invalid card number: %s", err.Error())
		}
		request.UserID = card.UserId
	}
	if request.UserID != 0 {
		db = db.Where("user_id = ?", request.UserID)
	}

	if err := db.Count(&cnt).Error; err != nil {
		return nil, 0, err
	}
	list := []models.Withdrawal{}
	if err := db.Scopes(models.ScopePaginate(request.Page, request.PageSize)).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	for _, ll := range list {
		resp = append(resp, response.WithdrawalRecordListResponseData{
			Id:         ll.Id,
			UserId:     ll.UserId,
			MerchantId: ll.MerchantId,
			Amount:     ll.Amount,
			Fee:        ll.Fee,
			Balance:    ll.Balance,
			CreatedAt:  ll.CreatedAt,
			UpdatedAt:  ll.UpdatedAt,
		})
	}
	return resp, cnt, nil
}

//如果是用户提现不需要扣除手续费
//如果是商户提现需要扣除手续费

func (withdrawal *WithdrawalService) Withdrawal(request request.WithdrawalRequest) (err error) {
	if request.UserId == 0 && request.MerchantId == 0 {
		return errors.New("user_id and merchant_id is required")
	}
	if request.Amount.LessThan(decimal.NewFromInt(0)) {
		return errors.New("amount is less than 0")
	}
	if request.MerchantId != 0 {
		merchantModel := models.Merchant{}
		if err := merchantModel.GetByID(request.MerchantId); err != nil {
			app.ZapLog.Error("merchant_id is invalid", zap.Error(err))
			return errors.New("merchant_id is invalid")
		}
		merchantBalanceData, err := merchantModel.GetMerchantBalance(request.MerchantId)
		if err != nil {
			app.ZapLog.Error("merchant_id is invalid", zap.Error(err))
			return errors.New("merchant_id is invalid")
		}

		if merchantBalanceData.CanWithdrawalBalance.LessThan(request.Amount) {
			app.ZapLog.Error("merchant can withdrawal balance is not enough", zap.Error(err))
			return errors.New("merchant can withdrawal balance is not enough")
		}

		//系统设置的提现手续费
		withdrawalFee := merchantBalanceData.FeeRate
		//实际对应的余额 = 申请金额/(1-手续费比例)
		actualAmount := request.Amount.Div(decimal.NewFromInt(1).Sub(withdrawalFee))
		//提现手续费
		feeAmount := actualAmount.Mul(withdrawalFee)

		withdrawalModel := models.Withdrawal{
			MerchantId: request.MerchantId,
			Amount:     request.Amount,                                               //申请提现金额
			Fee:        feeAmount,                                                    //手续费金额
			Balance:    merchantBalanceData.CanWithdrawalBalance.Sub(request.Amount), //剩余余额
		}

		if err := withdrawalModel.Insert(); err != nil {
			app.ZapLog.Error("create withdrawal record failed", zap.Error(err))
			return err
		}

		return nil
	}

	if request.UserId != 0 {
		userModel := models.User{}
		if err := userModel.GetById(request.UserId); err != nil {
			app.ZapLog.Error("user not found", zap.Error(err))
			return errors.New("user not found")
		}
		if userModel.Status == models.UserStatusDestroy {
			return errors.New("user is destroy")
		}

		if userModel.BindType == models.UserBindTypePassword {
			if request.Password == "" {
				return errors.New("password is empty")
			}
			if err := userModel.ComparePassword(request.Password); err != nil {
				app.ZapLog.Error("compare password failed", zap.Error(err))
				return errors.New("password is invalid")
			}
		} else if userModel.BindType == models.UserBindTypePhone {
			if userModel.SetWithdrawalCodeAt.Before(time.Now()) {
				return errors.New("code is expired, please set code again")
			}
			if userModel.WithdrawalCode != request.Code {
				return errors.New("code is invalid")
			}
		}

		userBalance, err := userModel.GetUserBalance(request.UserId)
		if err != nil {
			app.ZapLog.Error("user balance is invalid", zap.Error(err))
			return errors.New("user balance is invalid")
		}
		if userBalance.LessThan(request.Amount) {
			app.ZapLog.Error("user balance is not enough", zap.Error(err))
			return errors.New("user balance is not enough")
		}
		withdrawalModel := models.Withdrawal{
			UserId:  request.UserId,
			Amount:  request.Amount,
			Fee:     decimal.NewFromInt(0),
			Balance: userBalance.Sub(request.Amount),
		}
		if err := withdrawalModel.Insert(); err != nil {
			app.ZapLog.Error("create withdrawal record failed", zap.Error(err))
			return err
		}
		if userModel.BindType == models.UserBindTypePhone {
			userModel.WithdrawalCode = ""
			if err := userModel.Save(); err != nil {
				app.ZapLog.Error("save withdrawal record failed", zap.Error(err))
				return err
			}
		}

		return nil
	}
	return errors.New("user_id and merchant_id is required")
}

func (withdrawal *WithdrawalService) WithdrawalSendCode(request request.WithdrawalSendCodeRequest) (resp response.WithdrawalSendCodeResponse, err error) {
	userModel := models.User{}
	if err := userModel.GetById(request.UserId); err != nil {
		app.ZapLog.Error("user not found", zap.Error(err))
		return response.WithdrawalSendCodeResponse{}, errors.New("user not found")
	}
	// 生成4位数字验证码 (1000-9999)
	code := utils.RandInt(9000) + 1000
	userModel.WithdrawalCode = cast.ToString(code)
	//有效期 10 分钟
	userModel.SetWithdrawalCodeAt = time.Now().Add(10 * time.Minute)
	if err := userModel.Save(); err != nil {
		app.ZapLog.Error("save user failed", zap.Error(err))
		return response.WithdrawalSendCodeResponse{}, errors.New("save user failed")
	}
	resp.Url = fmt.Sprintf("https://wa.me/%s?text=%s", cast.ToString(userModel.PhoneNumber), cast.ToString(code))
	return resp, nil
}
