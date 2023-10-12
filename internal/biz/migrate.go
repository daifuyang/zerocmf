package biz

import (
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	err := new(User).AutoMigrate(db)
	if err != nil {
		return err
	}
	return nil
}
