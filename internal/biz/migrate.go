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

	err = new(SysDept).AutoMigrate(db)
	if err != nil {
		return err
	}

	err = new(SysMenu).AutoMigrate(db)
	if err != nil {
		return err
	}

	err = new(SysRole).AutoMigrate(db)
	if err != nil {
		return err
	}

	err = new(SysPost).AutoMigrate(db)
	if err != nil {
		return err
	}

	return nil
}
