package data

import (
	"context"
	"zerocmf/internal/biz"
)

type UserRepo struct {
	data *Data
}

// CreateUser implements biz.UserRepo.
func (*UserRepo) CreateUser(ctx context.Context, user *biz.User) (*biz.User, error) {
	panic("unimplemented")
}

func NewUserRepo(data *Data) biz.UserRepo {
	return &UserRepo{
		data: data,
	}
}
