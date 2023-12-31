package data

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"zerocmf/internal/biz"
	"zerocmf/internal/utils"

	v4Oauth2 "github.com/go-oauth2/oauth2/v4"
	v4Server "github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-redis/redis/v8"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

var (
	userCachePrefix = "cache:user:id:"
)

type userRepo struct {
	data *Data
}

// Find implements biz.UserRepo.
func (repo *userRepo) Find(ctx context.Context, listQuery *biz.UserListQuery) (*biz.Paginate, error) {

	var total int64 = 0

	userType := listQuery.UserType

	query := "user_type = ?"
	queryArgs := []interface{}{userType}

	tx := repo.data.db.Where(query, queryArgs...).Model(&biz.User{}).Count(&total)
	if tx.Error != nil {
		return nil, tx.Error
	}
	current := listQuery.Current
	pageSize := listQuery.PageSize

	offset := (current - 1) * pageSize
	var userList []*biz.User
	tx = repo.data.db.Where(query, queryArgs...).Offset(offset).Limit(pageSize).Find(&userList)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &biz.Paginate{
		Current:  current,
		PageSize: pageSize,
		Total:    total,
		Data:     userList,
	}, nil
}

// 验证token
func (repo *userRepo) ValidationBearerToken(req *http.Request) (v4Oauth2.TokenInfo, error) {
	srv := repo.data.srv
	return srv.ValidationBearerToken(req)
}

// 获取token
func (repo *userRepo) Token(ctx context.Context, user *biz.User) (*oauth2.Token, error) {

	config := repo.data.config.Oauth2

	srv := repo.data.srv
	oauth2Conf := repo.data.oauth2Conf

	authorizeRequest := &v4Server.AuthorizeRequest{
		ResponseType: "code",
		ClientID:     config.ClientID,
		RedirectURI:  config.RedirectURL,
		Scope:        "all",
		UserID:       strconv.FormatInt(int64(user.UserID), 10),
		// 其他参数根据需求设置
	}

	// 调用 GetAuthorizeToken 方法处理授权请求
	tokenInfo, err := srv.GetAuthorizeToken(ctx, authorizeRequest)
	if err != nil {
		return nil, err
	}

	code := tokenInfo.GetCode()
	token, err := oauth2Conf.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	return token, nil
}

// FindUserByAccount implements biz.UserRepo.
func (repo *userRepo) FindUserByAccount(ctx context.Context, account string) (*biz.User, error) {
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
func (repo *userRepo) FindUserByPhoneNumber(ctx context.Context, phoneNumber string) (*biz.User, error) {
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
func (repo *userRepo) FindOne(ctx context.Context, UserID int64) (*biz.User, error) {

	var user *biz.User
	key := fmt.Sprintf("%s%v", userCachePrefix, UserID)

	err := repo.data.RGet(ctx, &user, key)
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
func (repo *userRepo) CreateUser(ctx context.Context, user *biz.User) error {
	tx := repo.data.db.Create(&user)
	if tx.Error != nil {
		return tx.Error
	}

	key := fmt.Sprintf("%s%v", userCachePrefix, user.UserID)
	repo.data.RSet(ctx, key, user)
	return nil
}

func NewUserRepo(data *Data) biz.UserRepo {
	return &userRepo{
		data: data,
	}
}
