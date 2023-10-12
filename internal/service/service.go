package service

import (
	"zerocmf/configs"
	"zerocmf/internal/biz"
	"zerocmf/pkg/sms"

	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewContext)

// 用来组装上下文依赖
type Context struct {
	Config *configs.Config
	uc     *biz.Userusecase
	smsc   *biz.Smsusecase
	sms    sms.Sms
}

func NewContext(config *configs.Config, uc *biz.Userusecase, smsc *biz.Smsusecase) *Context {
	// 初始化短信
	client := sms.NewSms(config.Sms.AccessKeyId, config.Sms.AccessKeySecret, config.Sms.SignName, config.Sms.TemplateCode, config.Sms.Provider)
	return &Context{
		uc:   uc,
		smsc: smsc,
		sms:  client,
	}
}
