package services

import (
	"fmt"
	"time"

	"quickBillController/app"
	"quickBillController/config"
	"quickBillController/internal/merchant/params/request"
	"quickBillController/internal/merchant/params/response"
	"quickBillController/models"
	"quickBillController/utils/pcontext"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type MerchantService struct {
	PContext *pcontext.PMerchantContext
}

func NewMerchantServiceWithContext(c *gin.Context) *MerchantService {
	return &MerchantService{
		PContext: pcontext.ParseMerchantContext(c),
	}
}

func NewMerchantService() *MerchantService {
	return &MerchantService{}
}

func (merchant *MerchantService) MerchantUpdatePassword(request request.MerchantUpdatePasswordRequest) (err error) {
	merchantModel := models.Merchant{}
	if err := merchantModel.GetByID(merchant.PContext.MerchantId); err != nil {
		return err
	}
	if err := merchantModel.ComparePassword(request.OriginalPassword); err != nil {
		return fmt.Errorf("original password is incorrect")
	}
	pwd, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	merchantModel.Password = string(pwd)
	if err := merchantModel.Save(); err != nil {
		return err
	}
	return nil
}

func (merchant *MerchantService) RefreshToken() (resp *response.MerchantLoginResponse, err error) {
	merchantModel := models.Merchant{}
	if err := merchantModel.GetByID(merchant.PContext.MerchantId); err != nil {
		return nil, err
	}
	token, expiresAt, err := merchant.GenerateJWT(&merchantModel)
	if err != nil {
		return nil, err
	}
	return &response.MerchantLoginResponse{
		Id:        merchantModel.Id,
		UserName:  merchantModel.UserName,
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}
func (merchant *MerchantService) MerchantLogin(request request.MerchantLoginRequest) (resp *response.MerchantLoginResponse, err error) {
	merchantModel := models.Merchant{}
	if err := merchantModel.GetByUserName(request.UserName); err != nil {
		return nil, err
	}
	if err := merchantModel.ComparePassword(request.Password); err != nil {
		return nil, err
	}
	token, expiresAt, err := merchant.GenerateJWT(&merchantModel)
	if err != nil {
		return nil, err
	}
	return &response.MerchantLoginResponse{
		Id:        merchantModel.Id,
		UserName:  merchantModel.UserName,
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

func (merchant *MerchantService) GenerateJWT(merchantModel *models.Merchant) (tokenStr string, expiresAt int64, err error) {
	claims := MerchantJWTClaims{
		MerchantId: merchantModel.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.GetCfg().Server.JWTExpireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	expiresAt = claims.ExpiresAt.Time.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err = token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return "", 0, err
	}
	return tokenStr, expiresAt, nil
}

type MerchantJWTClaims struct {
	MerchantId int64 `json:"merchant_id"`
	jwt.RegisteredClaims
}

func (merchant *MerchantService) MerchantValidateJWT(tokenString string) (resp *MerchantJWTClaims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &MerchantJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token : %v", err)
	}

	if claims, ok := token.Claims.(*MerchantJWTClaims); ok && token.Valid {
		merchantModel := models.Merchant{}
		if err := merchantModel.GetByID(claims.MerchantId); err != nil {
			return nil, fmt.Errorf("invalid token")
		}
		if merchantModel.Id == 0 {
			return nil, fmt.Errorf("invalid token")
		}
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
func (merchant *MerchantService) MerchantUndoBillRecordList(request request.MerchantUndoBillRecordListRequest) (resp []response.MerchantBillRecordListResponseData, cnt int64, err error) {
	billModel := models.UndoBill{}
	db := app.DB.Model(&billModel).Where("merchant_id = ?", merchant.PContext.MerchantId).Order("id DESC")
	if request.OrderNumber != "" {
		db = db.Where("order_number ilike ?", "%"+request.OrderNumber+"%")
	}
	if request.StartAt != 0 && request.EndAt != 0 {
		db = db.Where("created_at >= ? AND created_at < ?", time.Unix(request.StartAt, 0), time.Unix(request.EndAt, 0))
	}
	if err := db.Count(&cnt).Error; err != nil {
		return nil, 0, err
	}
	list := []models.UndoBill{}
	if err := db.Scopes(models.ScopePaginate(request.Page, request.PageSize)).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	for _, bill := range list {
		resp = append(resp, response.MerchantBillRecordListResponseData{
			Id:           bill.Id,
			OrderNumber:  bill.OrderNumber,
			UserId:       bill.UserId,
			MerchantId:   bill.MerchantId,
			AliasNumber:  bill.AliasNumber,
			Amount:       bill.Amount,
			Balance:      bill.Balance,
			AfterBalance: bill.AfterBalance,
			CreatedAt:    bill.CreatedAt,
		})
	}
	return resp, cnt, nil
}
func (merchant *MerchantService) MerchantBillRecordList(request request.MerchantBillRecordListRequest) (resp []response.MerchantBillRecordListResponseData, cnt int64, err error) {
	billModel := models.Bill{}
	db := app.DB.Model(&billModel).Where("merchant_id = ?", merchant.PContext.MerchantId).Order("id DESC")
	if request.OrderNumber != "" {
		db = db.Where("order_number ilike ?", "%"+request.OrderNumber+"%")
	}
	if request.StartAt != 0 && request.EndAt != 0 {
		db = db.Where("created_at >= ? AND created_at < ?", time.Unix(request.StartAt, 0), time.Unix(request.EndAt, 0))
	}
	if err := db.Count(&cnt).Error; err != nil {
		return nil, 0, err
	}
	list := []models.Bill{}
	if err := db.Scopes(models.ScopePaginate(request.Page, request.PageSize)).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	for _, bill := range list {
		resp = append(resp, response.MerchantBillRecordListResponseData{
			Id:           bill.Id,
			OrderNumber:  bill.OrderNumber,
			UserId:       bill.UserId,
			MerchantId:   bill.MerchantId,
			AliasNumber:  bill.AliasNumber,
			Amount:       bill.Amount,
			Balance:      bill.Balance,
			AfterBalance: bill.AfterBalance,
			CreatedAt:    bill.CreatedAt,
		})
	}
	return resp, cnt, nil
}

func (merchant *MerchantService) MerchantSummary() (resp *response.MerchantSummaryResponse, err error) {
	merchantModel := models.Merchant{}
	merchantBalanceData, err := merchantModel.GetMerchantBalance(merchant.PContext.MerchantId)
	if err != nil {
		return nil, err
	}
	withdrawalModel := models.Withdrawal{}
	totalWithdrawalAmount, err := withdrawalModel.GetTotalWithdrawalAmountByMerchantID(merchant.PContext.MerchantId)
	if err != nil {
		return nil, err
	}
	return &response.MerchantSummaryResponse{
		FeeRate:                      merchantBalanceData.FeeRate,
		FrozenAmount:                 merchantBalanceData.FrozenAmount.RoundDown(2),
		OriginalFrozenAmount:         merchantBalanceData.OriginalFrozenAmount.RoundDown(2),
		CanWithdrawalBalance:         merchantBalanceData.CanWithdrawalBalance.RoundDown(2),
		OriginalCanWithdrawalBalance: merchantBalanceData.OriginalCanWithdrawalBalance.RoundDown(2),
		TotalAmount:                  merchantBalanceData.TotalAmount.RoundDown(2),
		TotalWithdrawalAmount:        totalWithdrawalAmount.RoundDown(2),
	}, nil
}

func (merchant *MerchantService) MerchantConsume(request request.MerchantConsumeRequest) (resp *response.MerchantConsumeResponse, err error) {
	if request.Amount.LessThan(decimal.NewFromInt(0)) {
		return nil, fmt.Errorf("amount is less than 0")
	}
	cardModel := models.Card{}
	if err := cardModel.GetByCardNumber(request.CardNumber); err != nil {
		return nil, err
	}
	userModel := models.User{}
	if err := userModel.GetById(cardModel.UserId); err != nil {
		return nil, err
	}
	if userModel.Status == models.UserStatusDestroy {
		return nil, fmt.Errorf("this card is destroyed")
	}
	userBalance, err := userModel.GetUserBalance(cardModel.UserId)
	if err != nil {
		return nil, err
	}
	if userBalance.LessThan(request.Amount) {
		return nil, fmt.Errorf("insufficient balance")
	}

	merchantModel := models.Merchant{}
	if err := merchantModel.GetByID(merchant.PContext.MerchantId); err != nil {
		return nil, err
	}

	billModel := models.Bill{}
	orderNumber := billModel.GenerateOrderNumber(userModel.Id, merchant.PContext.MerchantId)
	billModel = models.Bill{
		OrderNumber:  orderNumber,
		UserId:       userModel.Id,
		MerchantId:   merchant.PContext.MerchantId,
		CardNumber:   request.CardNumber,
		AliasNumber:  cardModel.AliasNumber,
		Amount:       request.Amount,
		Balance:      userBalance,
		AfterBalance: userBalance.Sub(request.Amount),
		CreatedAt:    time.Now(),
	}

	if err = billModel.Insert(); err != nil {
		app.ZapLog.Error("bill insert failed", zap.Error(err))
		return nil, fmt.Errorf("bill insert failed")
	}
	return &response.MerchantConsumeResponse{
		Id:           billModel.Id,
		OrderNumber:  billModel.OrderNumber,
		UserId:       billModel.UserId,
		MerchantId:   billModel.MerchantId,
		AliasNumber:  billModel.AliasNumber,
		Amount:       billModel.Amount,
		Balance:      billModel.Balance,
		AfterBalance: billModel.AfterBalance,
		CreatedAt:    billModel.CreatedAt,
	}, nil
}

func (merchant *MerchantService) MerchantUndoConsume(request request.MerchantUndoConsumeRequest) (resp models.UndoBill, err error) {
	billModel := models.Bill{}
	if err := billModel.GetByOrderNumber(request.OrderNumber); err != nil {
		return resp, err
	}
	if billModel.Id == 0 {
		return resp, fmt.Errorf("bill not found")
	}
	if billModel.MerchantId != merchant.PContext.MerchantId {
		return resp, fmt.Errorf("bill merchant id is not equal to merchant")
	}

	//只有冻结时间内的订单可撤销
	frozenHours := config.GetCfg().Merchant.FrozenHours
	if billModel.CreatedAt.Before(time.Now().Add(-time.Duration(frozenHours) * time.Hour)) {
		return resp, fmt.Errorf("only one day order can be undone")
	}

	userModel := models.User{}
	if err := userModel.GetById(billModel.UserId); err != nil {
		return resp, err
	}
	if userModel.Status == models.UserStatusDestroy {
		return resp, fmt.Errorf("this card is destroyed")
	}
	userModelBalance, err := userModel.GetUserBalance(billModel.UserId)
	if err != nil {
		return resp, err
	}

	undoBillModel := models.UndoBill{
		OrderNumber:  request.OrderNumber,
		UserId:       billModel.UserId,
		MerchantId:   billModel.MerchantId,
		AliasNumber:  billModel.AliasNumber,
		CardNumber:   billModel.CardNumber,
		Amount:       billModel.Amount,
		Balance:      userModelBalance,
		AfterBalance: userModelBalance.Add(billModel.Amount),
		CreatedAt:    billModel.CreatedAt,
	}

	err = app.DB.Transaction(func(tx *gorm.DB) error {
		//记录撤销消费记录
		if err := tx.Create(&undoBillModel).Error; err != nil {
			app.ZapLog.Error("create undo bill failed", zap.Error(err))
			return fmt.Errorf("create undo bill failed")
		}

		//删除消费记录
		if err := tx.Delete(&billModel).Error; err != nil {
			app.ZapLog.Error("delete bill failed", zap.Error(err))
			return fmt.Errorf("delete bill failed")
		}
		return nil
	})
	return undoBillModel, err
}

func (merchant *MerchantService) MerchantQueryCardBalance(request request.MerchantQueryCardBalanceRequest) (resp *response.MerchantQueryCardBalanceResponse, err error) {
	cardModel := models.Card{}
	if err := cardModel.GetByCardNumber(request.CardNumber); err != nil {
		return nil, err
	}
	userModel := models.User{}
	if err := userModel.GetById(cardModel.UserId); err != nil {
		return nil, err
	}
	if userModel.Status == models.UserStatusDestroy {
		return nil, fmt.Errorf("this card is destroyed")
	}
	userBalance, err := userModel.GetUserBalance(cardModel.UserId)
	if err != nil {
		return nil, err
	}
	resp = &response.MerchantQueryCardBalanceResponse{
		Balance: userBalance,
	}
	return resp, nil
}
