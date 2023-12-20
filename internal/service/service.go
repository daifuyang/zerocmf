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
	Config *configs.Config
	Logger *zap.Logger
	uc     *biz.Userusecase
	smsc   *biz.Smsusecase
	dc     *biz.Depatmentusecase
	rc     *biz.Roleusecase
	mc     *biz.Menusecase
	postuc *biz.Postusecase
}

func NewContext(config *configs.Config, uc *biz.Userusecase, smsc *biz.Smsusecase, dc *biz.Depatmentusecase, mc *biz.Menusecase, rc *biz.Roleusecase, postuc *biz.Postusecase) *Context {

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	// 初始化短信
	return &Context{
		Config: config,
		Logger: logger,
		uc:     uc,
		smsc:   smsc,
		dc:     dc,
		mc:     mc,
		rc:     rc,
		postuc: postuc,
	}
}
