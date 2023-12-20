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
			user.GET("/current_user", svcCtx.CurrentUser)

			// 角色管理
			user.GET("/roles", service.NewRole(svcCtx).List)
			user.GET("/roles/:id", service.NewRole(svcCtx).Show)
			user.POST("/roles", service.NewRole(svcCtx).Add)
			user.POST("/roles/:id", service.NewRole(svcCtx).Update)
			user.DELETE("/roles/:id", service.NewRole(svcCtx).Delete)

			// 部门管理
			user.GET("/post", service.NewPost(svcCtx).List)
			user.GET("/post/:id", service.NewPost(svcCtx).Show)
			user.POST("/post", service.NewPost(svcCtx).Add)
			user.POST("/post/:id", service.NewPost(svcCtx).Update)
			user.DELETE("/post/:id", service.NewPost(svcCtx).Delete)

			// 管理员相关
			user.GET("/admins", service.NewAdmin(svcCtx).List)
			user.GET("/admins/:id", service.NewAdmin(svcCtx).List)
			user.POST("/admins", service.NewAdmin(svcCtx).Add)
			user.POST("/admins/:id", service.NewAdmin(svcCtx).Update)
			user.DELETE("/admins/:id", service.NewAdmin(svcCtx).Delete)
		}

		// 系统相关
		system := v1.Group("/system")
		{
			system.Use(svcCtx.AuthMiddleware)

			// 菜单相关
			system.GET("/menus", service.NewMenu(svcCtx).Tree)
			system.GET("/menus/:id", service.NewMenu(svcCtx).Show)
			system.POST("/menus", service.NewMenu(svcCtx).Add)
			system.POST("/menus/:id", service.NewMenu(svcCtx).Update)
			system.DELETE("/menus/:id", service.NewMenu(svcCtx).Delete)

			// 部门相关
			system.GET("/departments", service.NewDeparment(svcCtx).Tree)
			system.GET("/departments/:id", service.NewDeparment(svcCtx).Show)
			system.POST("/departments", service.NewDeparment(svcCtx).Add)
			system.POST("/departments/:id", service.NewDeparment(svcCtx).Update)

			// 权限相关
			// system.GET("/permissions/:id", service.NewAuthz(svcCtx).Index)
			// system.POST("/permissions/:id", service.NewAuthz(svcCtx).Save)

			// 系统监控
			system.GET("/monitor/server", service.NewMonitor(svcCtx).Index)
		}
	}

}
