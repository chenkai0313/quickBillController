package models

import (
	"time"

	"quickBillController/app"
)

type Card struct {
	Id          int64     `gorm:"primarykey" json:"id"`
	Number      string    `gorm:"type:varchar(1000);not null;unique;default:''"  json:"number"`
	AliasNumber string    `gorm:"type:varchar(1000);not null;unique;default:''"  json:"alias_number"`
	UserId      int64     `gorm:"type:int8;not null;default:0;" json:"uid"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (m *Card) GetByID(id int64) error {
	if err := app.DB.Where("id = ?", id).First(m).Error; err != nil {
		return err
	}
	return nil
}

func (m *Card) Insert() error {
	return app.DB.Create(m).Error
}

func (m *Card) Save() error {
	return app.DB.Save(m).Error
}

func (m *Card) GetByCardNumber(cardNumber string) error {
	if err := app.DB.Where("number = ?", cardNumber).First(m).Error; err != nil {
		return err
	}
	return nil
}

func (m *Card) GetByCardNumberOrAliasNumber(cardNumber string) error {
	if err := app.DB.Where("number = ?", cardNumber).Or("alias_number=?", cardNumber).First(m).Error; err != nil {
		return err
	}
	return nil
}

func (m *Card) GetByAliasNumber(aliasNumber string) error {
	if err := app.DB.Where("alias_number = ?", aliasNumber).First(m).Error; err != nil {
		return err
	}
	return nil
}
