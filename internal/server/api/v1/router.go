package v1

import (
	"zerocmf/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterHTTPServer(router *gin.Engine, svcCtx *service.Context) {

	v1 := router.Group("/api/v1")
	{
		// 用户注册
		user := v1.Group("/user")
		{
			user.POST("/register", svcCtx.Register)
			user.GET("/send_register_code", svcCtx.SendRegisterCode)
		}
	}

}
