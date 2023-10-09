package svc

import (
	"zerocmf/configs"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Engine *gin.Engine
	Db     *gorm.DB
	Config configs.Config
}
