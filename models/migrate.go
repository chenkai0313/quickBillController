package models

import (
	"fmt"

	"quickBillController/app"
)

func Migrate() {
	err := app.DB.AutoMigrate(&Admin{}, &Bill{}, &Card{}, &Merchant{}, &User{}, &Withdrawal{}, &Topup{}, &UndoBill{})

	if err != nil {
		panic(fmt.Sprintf("Database model table migreate failed %v", err.Error()))
	}

	err = InitAdmin()
	if err != nil {
		panic(fmt.Sprintf("Admin initialization failed %v", err.Error()))
	}

}
