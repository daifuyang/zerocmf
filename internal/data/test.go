package data

import (
	"context"
	"zerocmf/internal/biz"
)

type testRepo struct {
	data *Data
}

// NewTestRepo .
func NewTestRepo(data *Data) biz.TestRepo {
	return &testRepo{
		data: data,
	}
}

// ListArticle implements biz.TestRepo.
func (*testRepo) ListArticle(ctx context.Context) ([]*biz.Test, error) {
	return []*biz.Test{}, nil
}
