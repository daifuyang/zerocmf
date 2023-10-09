package service

import (
	"context"
	"fmt"
	"zerocmf/internal/biz"
)

type TestService struct {
	test *biz.TestUsecase
}

func NewTestService(test *biz.TestUsecase) *TestService {
	return &TestService{
		test: test,
	}
}

func (s *TestService) CreateArticle(ctx context.Context) error {
	fmt.Println("ssss")
	s.test.List(ctx)
	return nil
}
