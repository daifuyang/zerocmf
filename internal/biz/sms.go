package biz

import (
	"context"
	"fmt"
	"zerocmf/configs"
	"zerocmf/pkg/sms"
)

type Sms struct {
	Ctx       context.Context
	Code      uint
	Account   string
	ClientIP  string
	Type      string
	SceneCode string
}

// 定义data层接口
type SmsRepo interface {
	CacheCode(sms *Sms) error  // 缓存短信验证码
	VerifyCode(sms *Sms) error // 验证短信验证码
}

type Smsusecase struct {
	repo SmsRepo
	sms  sms.Sms
}

func NewSmsUsecase(repo SmsRepo, config *configs.Config) *Smsusecase {
	client := sms.NewSms(config.Sms.AccessKeyId, config.Sms.AccessKeySecret, config.Sms.SignName, config.Sms.TemplateCode, config.Sms.Provider)
	return &Smsusecase{repo: repo, sms: client}
}

// 发送验证码

func (s *Smsusecase) SendSms(smsCtx *Sms) error {
	err := s.repo.CacheCode(smsCtx)
	if err != nil {
		return err
	}

	templateParam := fmt.Sprintf("{\"code\":\"%d\"}", smsCtx.Code)
	return s.sms.SendSms(smsCtx.Account, templateParam)
}

// 验证验证码
func (s *Smsusecase) VerifyCode(smsCtx *Sms) error {
	err := s.repo.VerifyCode(smsCtx)
	return err
}
