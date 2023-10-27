package service

import (
	"zerocmf/configs"
	"zerocmf/internal/biz"
	"zerocmf/pkg/casbinz"
	zOauth "zerocmf/pkg/oauth2"

	"github.com/casbin/casbin/v2"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/google/wire"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewContext)

// 用来组装上下文依赖
type Context struct {
	Config      *configs.Config
	Logger      *zap.Logger
	e           *casbin.Enforcer
	uc          *biz.Userusecase
	smsc        *biz.Smsusecase
	dc          *biz.Depatmentusecase
	mc          *biz.Menusecase
	srv         *server.Server
	oauthConfig oauth2.Config
}

func NewContext(config *configs.Config, uc *biz.Userusecase, smsc *biz.Smsusecase, dc *biz.Depatmentusecase, mc *biz.Menusecase) *Context {

	oauthConfig, srv := zOauth.NewServer(config)

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	e := casbinz.NewAdapter(config)

	// 初始化短信
	return &Context{
		Config:      config,
		Logger:      logger,
		uc:          uc,
		e:           e,
		smsc:        smsc,
		dc:          dc,
		mc:          mc,
		srv:         srv,
		oauthConfig: oauthConfig,
	}
}
