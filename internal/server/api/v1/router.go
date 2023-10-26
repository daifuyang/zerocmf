package v1

import (
	"zerocmf/internal/service"

	"github.com/gin-gonic/gin"
	ginserver "github.com/go-oauth2/gin-server"
)

func RegisterHTTPServer(router *gin.Engine, svcCtx *service.Context) {

	// srv.SetUserAuthorizationHandler(func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	// 	return "1", nil
	// })

	auth := router.Group("/oauth2")
	{
		// auth.GET("/authorize", ginserver.HandleAuthorizeRequest)
		// auth.GET("/callback", func(c *gin.Context) {
		// 	svcCtx.Callback(c, oauthConfig)
		// })
		auth.POST("/token", ginserver.HandleTokenRequest)
	}

	v1 := router.Group("/api/v1")
	{
		// 用户注册
		user := v1.Group("/user")
		{
			user.POST("/register", svcCtx.Register)                  // 注册
			user.GET("/send_register_code", svcCtx.SendRegisterCode) // 发送注册验证码
			user.POST("/login", svcCtx.Login)                        // 登录

			// 验证用户身份信息

			user.Use(svcCtx.AuthMiddleware)
			user.GET("/current_user", func(ctx *gin.Context) {

			})
		}

		system := v1.Group("/system")
		{
			system.Use(svcCtx.AuthMiddleware)

			// 菜单相关
			system.GET("/menus", service.NewMenu(svcCtx).Tree)
			system.GET("/menus/:id", service.NewMenu(svcCtx).Show)
			system.POST("/menus", service.NewMenu(svcCtx).Add)
			system.POST("/menus/:id", service.NewMenu(svcCtx).Update)

			// 部门相关
			system.GET("/departments", service.NewDeparment(svcCtx).Tree)
			system.GET("/departments/:id", service.NewDeparment(svcCtx).Show)
			system.POST("/departments", service.NewDeparment(svcCtx).Add)
			system.POST("/departments/:id", service.NewDeparment(svcCtx).Update)
		}
	}

}
