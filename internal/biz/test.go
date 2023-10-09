package biz

import "context"

type Test struct {
	ID   int64
	Name string
}

type TestRepo interface {
	// db
	ListArticle(ctx context.Context) ([]*Test, error)
}

type TestUsecase struct {
	repo TestRepo
}

func NewTestUsecase(repo TestRepo) *TestUsecase {
	return &TestUsecase{repo: repo}
}

func (uc *TestUsecase) List(ctx context.Context) (ps []*Test, err error) {
	ps, err = uc.repo.ListArticle(ctx)
	if err != nil {
		return
	}
	return
}
