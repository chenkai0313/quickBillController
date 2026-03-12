package models

import "gorm.io/gorm"

func ScopePaginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page < 1 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100 // Maximum limit 100
		case pageSize <= 0:
			pageSize = 10 // Default 10 items
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
