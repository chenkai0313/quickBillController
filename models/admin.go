package models

import (
	"time"

	"quickBillController/app"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Admin struct {
	Id        int64     `gorm:"primarykey" json:"id"`
	UserName  string    `gorm:"type:varchar(100);not null;unique;" json:"user_name"`
	Password  string    `json:"password" gorm:"type:varchar(1000);not null;"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Check if admin user exists in the table, if not, initialize default admin user with username 'admin' and password '1qaz2wsx'
func InitAdmin() error {
	var admin Admin
	if err := app.DB.Where("user_name = ?", "admin").First(&admin).Error; err != nil {
		admin.UserName = "admin"
		password, err := bcrypt.GenerateFromPassword([]byte("1qaz2wsx"), bcrypt.DefaultCost)
		if err != nil {
			app.ZapLog.Error("Admin password generation failed", zap.Error(err))
			return err
		}
		admin.Password = string(password)
		if err := app.DB.Create(&admin).Error; err != nil {
			app.ZapLog.Error("Admin creation failed", zap.Error(err))
			return err
		}
		app.ZapLog.Info("Admin initialized", zap.String("admin", admin.UserName))
		return nil
	}
	return nil
}

func (a *Admin) GetByUserName(userName string) error {
	if err := app.DB.Where("user_name = ?", userName).First(a).Error; err != nil {
		return err
	}
	return nil
}

func (a *Admin) GetById(id int64) error {
	if err := app.DB.Where("id = ?", id).First(a).Error; err != nil {
		return err
	}
	return nil
}
func (a *Admin) ComparePassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password)); err != nil {
		return err
	}
	return nil
}

func (a *Admin) Save() error {
	return app.DB.Save(a).Error
}
