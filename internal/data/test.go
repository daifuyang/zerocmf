package data

import (
	"context"
	"fmt"
	"zerocmf/internal/biz"
)

type testRepo struct {
	data *Data
}

// NewArticleRepo .
func NewTestRepo(data *Data) biz.TestRepo {
	return &testRepo{
		data: data,
	}
}

// ListArticle implements biz.TestRepo.
func (*testRepo) ListArticle(ctx context.Context) ([]*biz.Test, error) {
	fmt.Println("ListArticle")
	return []*biz.Test{}, nil
}
