package biz

import (
	"context"
)

type Sms struct {
	Ctx       context.Context
	Code      string
	Account   string
	Type      string
	SceneCode string
}

// 定义data层接口
type SmsRepo interface {
	CacheCode(sms *Sms) error // 缓存验证码
}

type Smsusecase struct {
	repo SmsRepo
}

func NewSmsUsecase(repo SmsRepo) *Smsusecase {
	return &Smsusecase{repo: repo}
}

func (s *Smsusecase) SendSms(smsCtx *Sms) error {
	err := s.repo.CacheCode(smsCtx)
	return err
}
