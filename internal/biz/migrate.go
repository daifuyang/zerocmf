package biz

import (
	"zerocmf/configs"

	"gorm.io/gorm"
)

// 数据库迁移
func AutoMigrate(db *gorm.DB, config *configs.Config) error {
	salt := config.Mysql.Salt
	err := new(User).AutoMigrate(db, salt)
	if err != nil {
		return err
	}

	err = new(Dept).AutoMigrate(db)
	if err != nil {
		return err
	}

	err = new(Menu).AutoMigrate(db)
	if err != nil {
		return err
	}

	err = new(Role).AutoMigrate(db)
	if err != nil {
		return err
	}

	err = new(Post).AutoMigrate(db)
	if err != nil {
		return err
	}

	err = new(UserRole).AutoMigrate(db)
	if err != nil {
		return err
	}

	err = new(UserPost).AutoMigrate(db)
	if err != nil {
		return err
	}

	return nil
}
