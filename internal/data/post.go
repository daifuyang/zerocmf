package data

import (
	"context"
	"zerocmf/internal/biz"
)

type PostRepo struct {
	data *Data
}

// Find implements biz.PostRepo.
func (*PostRepo) Find(ctx context.Context, listQuery *biz.SysPostListQuery) (interface{}, error) {
	panic("unimplemented")
}

// FindOne implements biz.PostRepo.
func (*PostRepo) FindOne(ctx context.Context, id int64) (*biz.SysPost, error) {
	panic("unimplemented")
}

// Insert implements biz.PostRepo.
func (*PostRepo) Insert(ctx context.Context, menu *biz.SysPost) (err error) {
	panic("unimplemented")
}

// Update implements biz.PostRepo.
func (*PostRepo) Update(ctx context.Context, menu *biz.SysPost) (err error) {
	panic("unimplemented")
}

// Delete implements biz.PostRepo.
func (*PostRepo) Delete(ctx context.Context, id int64) error {
	panic("unimplemented")
}

func NewPostRepo(data *Data) biz.PostRepo {
	return &PostRepo{data: data}
}
