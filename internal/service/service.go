package service

import (
	"zerocmf/configs"
	"zerocmf/internal/biz"

	"github.com/google/wire"
	"go.uber.org/zap"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewContext)

// 用来组装上下文依赖
type Context struct {
	Config  *configs.Config
	Logger  *zap.Logger
	useruc  *biz.Userusecase
	smsc    *biz.Smsusecase
	deptuc  *biz.Depatmentusecase
	roleuc  *biz.Roleusecase
	menusuc *biz.Menusecase
	postuc  *biz.Postusecase
}

func NewContext(config *configs.Config, useruc *biz.Userusecase, smsc *biz.Smsusecase, deptuc *biz.Depatmentusecase, menusuc *biz.Menusecase, roleuc *biz.Roleusecase, postuc *biz.Postusecase) *Context {

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	// 初始化短信
	return &Context{
		Config:  config,
		Logger:  logger,
		useruc:  useruc,
		smsc:    smsc,
		deptuc:  deptuc,
		menusuc: menusuc,
		roleuc:  roleuc,
		postuc:  postuc,
	}
}
