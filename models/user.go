package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"quickBillController/app"

	"github.com/shopspring/decimal"
)

type User struct {
	Id                  int64     `gorm:"primarykey" json:"id"`
	CardNumber          string    `json:"card_number" gorm:"type:varchar(1000);not null;index;default:''" `
	AliasNumber         string    `json:"alias_number" gorm:"type:varchar(1000);not null;index;default:''" `
	Status              int       `json:"status" gorm:"type:int8;not null;default:0;index;comment:'0:destroyed 1:using'"`
	Password            string    `json:"password" gorm:"type:varchar(100);not null;default:''" `
	PhoneNumber         string    `json:"phone_number" gorm:"type:varchar(1000);not null;default:''" `
	BindType            int       `json:"bind_type" gorm:"type:int8;not null;default:0;index;comment:'0:null 1:password 2:phone'"`
	WithdrawalCode      string    `json:"withdrawal_code" gorm:"type:varchar(100);not null;default:''" `
	SetWithdrawalCodeAt time.Time `json:"set_withdrawal_code_at"`
	DestroyAt           time.Time `json:"destroy_at"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

const UserBindTypeNull = 0     //null
const UserBindTypePassword = 1 //password
const UserBindTypePhone = 2    //phone

func (m *User) ComparePassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(password)); err != nil {
		return err
	}
	return nil
}

const UserStatusUsing = 1   //using
const UserStatusDestroy = 0 //destroyed

func (m *User) GetById(id int64) error {
	if err := app.DB.Where("id = ?", id).First(m).Error; err != nil {
		return err
	}
	return nil
}

func (m *User) GetByCardNumber(cardNumber string) error {
	if err := app.DB.Where("card_number = ?", cardNumber).First(m).Error; err != nil {
		return err
	}
	return nil
}
func (m *User) Insert() error {
	return app.DB.Create(m).Error
}
func (m *User) Save() error {
	return app.DB.Save(m).Error
}

// update user status to destroy if card number exists and status is using
func (m *User) UpdateUsingStatusByCardNumber(cardNumber string) error {
	if err := app.DB.Model(&m).Where("card_number = ? and status = ?", cardNumber, UserStatusUsing).Updates(map[string]interface{}{"status": UserStatusDestroy, "destroy_at": time.Now()}).Error; err != nil {
		return err
	}
	return nil
}

// balance = 用户所有充值-所有消费-所有用户提现
func (m *User) GetUserBalance(id int64) (balance decimal.Decimal, err error) {
	// 1. 查询用户所有充值金额总和
	var topupAmount struct {
		Total decimal.Decimal
	}
	if err := app.DB.Model(&Topup{}).
		Where("user_id = ?", id).
		Select("COALESCE(SUM(amount), 0) as total").
		Scan(&topupAmount).Error; err != nil {
		return decimal.NewFromInt(0), err
	}

	// 2. 查询用户所有消费金额总和
	var billAmount struct {
		Total decimal.Decimal
	}
	if err := app.DB.Model(&Bill{}).
		Where("user_id = ?", id).
		Select("COALESCE(SUM(amount), 0) as total").
		Scan(&billAmount).Error; err != nil {
		return decimal.NewFromInt(0), err
	}

	// 3. 查询用户所有提现金额总和
	var withdrawalAmount struct {
		Total decimal.Decimal
	}
	if err := app.DB.Model(&Withdrawal{}).
		Where("user_id = ?", id).
		Select("COALESCE(SUM(amount), 0) as total").
		Scan(&withdrawalAmount).Error; err != nil {
		return decimal.NewFromInt(0), err
	}

	// 4. 计算余额：充值总和 - 消费总和 - 提现总和
	balance = topupAmount.Total.Sub(billAmount.Total).Sub(withdrawalAmount.Total)

	return balance, nil
}
