package data

import (
	"context"
	"errors"
	"fmt"
	"zerocmf/internal/biz"
	"zerocmf/internal/utils"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	userCachePrefix = "cache:user:userId:"
)

type UserRepo struct {
	data *Data
}

// FindUserByAccount implements biz.UserRepo.
func (repo *UserRepo) FindUserByAccount(ctx context.Context, account string) (*biz.User, error) {
	user := &biz.User{}
	query := "login_name = ?"
	queryArgs := []interface{}{account}

	if utils.AccountType(account) == "phone" {
		query = "phone_number = ?"
	} else if utils.AccountType(account) == "email" {
		query = "email = ?"
	}
	tx := repo.data.db.Where(query, queryArgs...).First(user)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("该账号不存在！")
		}
		return nil, tx.Error
	}
	return user, nil
}

// FindUserByPhoneNumber implements biz.UserRepo.
func (repo *UserRepo) FindUserByPhoneNumber(ctx context.Context, phoneNumber string) (*biz.User, error) {
	user := &biz.User{}
	tx := repo.data.db.Where("phone_number = ?", phoneNumber).First(user)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, tx.Error
	}
	return user, nil
}

// FindUserByUserID implements biz.UserRepo.
func (repo *UserRepo) FindOne(ctx context.Context, UserID uint64) (*biz.User, error) {

	var user *biz.User
	key := fmt.Sprintf("%s%v", userCachePrefix, UserID)

	err := repo.data.RGet(ctx, user, key)
	if err != redis.Nil {
		return user, nil
	}

	tx := repo.data.db.Where("user_id = ?", UserID).First(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}

	// 将结果存入redis
	err = repo.data.RSet(ctx, key, user)
	if err != nil {
		return nil, err
	}

	return user, nil

}

// CreateUser implements biz.UserRepo.
func (repo *UserRepo) CreateUser(ctx context.Context, user *biz.User) error {
	tx := repo.data.db.Create(&user)
	if tx.Error != nil {
		return tx.Error
	}

	key := fmt.Sprintf("%s%v", userCachePrefix, user.UserID)
	repo.data.RSet(ctx, key, user)
	return nil
}

func NewUserRepo(data *Data) biz.UserRepo {
	return &UserRepo{
		data: data,
	}
}
