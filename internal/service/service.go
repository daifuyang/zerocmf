package service

import (
	"zerocmf/configs"
	"zerocmf/internal/biz"
	zOauth "zerocmf/pkg/oauth2"

	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/google/wire"
	"golang.org/x/oauth2"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewContext)

// 用来组装上下文依赖
type Context struct {
	Config      *configs.Config
	uc          *biz.Userusecase
	smsc        *biz.Smsusecase
	dc          *biz.Depatmentusecase
	mc          *biz.Menusecase
	srv         *server.Server
	oauthConfig oauth2.Config
}

func NewContext(config *configs.Config, uc *biz.Userusecase, smsc *biz.Smsusecase, dc *biz.Depatmentusecase, mc *biz.Menusecase) *Context {

	oauthConfig, srv := zOauth.NewServer(config)

	// 初始化短信
	return &Context{
		Config:      config,
		uc:          uc,
		smsc:        smsc,
		dc:          dc,
		mc:          mc,
		srv:         srv,
		oauthConfig: oauthConfig,
	}
}
