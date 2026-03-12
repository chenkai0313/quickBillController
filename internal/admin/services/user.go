package services

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"quickBillController/app"
	"quickBillController/internal/admin/params/request"
	"quickBillController/internal/admin/params/response"
	"quickBillController/models"

	"go.uber.org/zap"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (user *UserService) UserList(request request.UserListRequest) (resp []response.UserListResponseData, cnt int64, err error) {
	userModel := models.User{}
	db := app.DB.Model(&userModel).Order("id DESC")
	if request.AliasNumber != "" {
		db = db.Where("alias_number iLIKE ?", "%"+request.AliasNumber+"%")
	}
	if request.PhoneNumber != "" {
		db = db.Where("phone_number iLIKE ?", "%"+request.PhoneNumber+"%")
	}
	var list []models.User
	if err := db.Count(&cnt).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Scopes(models.ScopePaginate(request.Page, request.PageSize)).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	for _, user := range list {
		mm := models.User{}
		ba, _ := mm.GetUserBalance(user.Id)
		resp = append(resp, response.UserListResponseData{
			Id:          user.Id,
			CardNumber:  user.CardNumber,
			AliasNumber: user.AliasNumber,
			Status:      user.Status,
			Balance:     ba,
			PhoneNumber: user.PhoneNumber,
			DestroyAt:   user.DestroyAt,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		})
	}
	return resp, cnt, nil
}

// query card if exists, if not create, if exists update card user_id data
// if card exists, create user and update card user_id data
func (user *UserService) UserCreate(request request.UserCreateRequest) (resp *response.UserCreateResponse, err error) {
	cardModel := models.Card{}
	_ = cardModel.GetByCardNumber(request.CardNumber)
	if cardModel.Id == 0 {
		cardModel.AliasNumber = request.AliasNumber
		cardModel.Number = request.CardNumber
		if err := cardModel.Insert(); err != nil {
			return nil, err
		}
	}
	// if card exists, update user status to destroy
	userModel := models.User{}
	if err := userModel.UpdateUsingStatusByCardNumber(request.CardNumber); err != nil {
		return nil, err
	}

	pwd := ""
	if request.BindType == models.UserBindTypePassword {
		if request.Password != "" {
			pwdByte, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
			if err != nil {
				return nil, err
			}
			pwd = string(pwdByte)
		}
	}
	userModel = models.User{
		CardNumber:  request.CardNumber,
		AliasNumber: request.AliasNumber,
		Password:    pwd,
		PhoneNumber: request.PhoneNumber,
		BindType:    request.BindType,
		Status:      models.UserStatusUsing,
	}
	if err := userModel.Insert(); err != nil {
		return nil, err
	}

	cardModel.UserId = userModel.Id
	if err := cardModel.Save(); err != nil {
		return nil, err
	}

	if !request.TopAmount.IsZero() {
		userBalance, err := userModel.GetUserBalance(cardModel.UserId)
		if err != nil {
			return nil, err
		}

		topupModel := models.Topup{
			UserId:    userModel.Id,
			Amount:    request.TopAmount,
			Balance:   userBalance,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := topupModel.Insert(); err != nil {
			return nil, err
		}
	}

	return &response.UserCreateResponse{
		Id:          userModel.Id,
		AliasNumber: cardModel.AliasNumber,
		CreatedAt:   userModel.CreatedAt,
		UpdatedAt:   userModel.UpdatedAt,
	}, nil
}

func (user *UserService) UserReadCards(request request.UserReadCardsRequest) (resp *response.UserReadCardsResponse, err error) {
	cardModel := models.Card{}
	if err := cardModel.GetByCardNumber(request.CardNumber); err != nil {
		return &response.UserReadCardsResponse{}, nil
	}
	if cardModel.Id == 0 {
		return &response.UserReadCardsResponse{}, nil
	}

	userModel := models.User{}
	if err := userModel.GetById(cardModel.UserId); err != nil {
		return &response.UserReadCardsResponse{}, nil
	}
	if userModel.Id == 0 {
		return &response.UserReadCardsResponse{}, nil
	}

	if userModel.Status == models.UserStatusDestroy {
		return &response.UserReadCardsResponse{}, nil
	}

	userBalance, err := userModel.GetUserBalance(cardModel.UserId)
	if err != nil {
		app.ZapLog.Error("GetUserBalance", zap.Error(err))
		return &response.UserReadCardsResponse{}, err
	}

	resp = &response.UserReadCardsResponse{
		Id:          userModel.Id,
		CardNumber:  cardModel.Number,
		BindType:    userModel.BindType,
		PhoneNumber: userModel.PhoneNumber,
		Status:      userModel.Status,
		DestroyAt:   userModel.DestroyAt,
		AliasNumber: cardModel.AliasNumber,
		Balance:     userBalance,
		CreatedAt:   userModel.CreatedAt,
		UpdatedAt:   userModel.UpdatedAt,
		BillList:    make([]response.UserBillListResponseData, 0),
	}

	billModel := models.Bill{}
	list, err := billModel.GetListByUserId(userModel.Id)
	if err != nil {
		return nil, err
	}
	for _, bill := range list {
		resp.BillList = append(resp.BillList, response.UserBillListResponseData{
			Id:           bill.Id,
			OrderNumber:  bill.OrderNumber,
			Amount:       bill.Amount,
			Balance:      bill.Balance,
			MerchantId:   bill.MerchantId,
			AfterBalance: bill.AfterBalance,
			CreatedAt:    bill.CreatedAt,
			MerchantName: bill.MerchantName,
		})
	}

	topupModel := models.Topup{}
	listTopup, err := topupModel.GetListByUserId(userModel.Id)
	if err != nil {
		return nil, err
	}
	for _, topup := range listTopup {
		resp.TopUpList = append(resp.TopUpList, response.UserTopUpListResponseData{
			Id:        topup.Id,
			Amount:    topup.Amount,
			Balance:   topup.Balance,
			CreatedAt: topup.CreatedAt,
		})
	}

	withdrawalModel := models.Withdrawal{}
	listWithdrawal, err := withdrawalModel.GetListByUserId(userModel.Id)
	if err != nil {
		return nil, err
	}
	for _, withdrawal := range listWithdrawal {
		resp.WithdrawalList = append(resp.WithdrawalList, response.UserWithdrawalListResponseData{
			Id:        withdrawal.Id,
			Amount:    withdrawal.Amount,
			Balance:   withdrawal.Balance,
			CreatedAt: withdrawal.CreatedAt,
		})
	}

	return resp, nil
}
