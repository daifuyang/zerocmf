package service

import (
	"strings"
	"zerocmf/internal/biz"
	"zerocmf/internal/utils"
	"zerocmf/pkg/hashed"
	"zerocmf/pkg/response"
	"zerocmf/pkg/sms"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// 用户注册
func (s *Context) Register(c *gin.Context) {
	ctx := c.Request.Context()

	req := new(biz.Register)
	c.ShouldBind(&req)

	account := req.Account
	if strings.TrimSpace(account) == "" {
		response.Error(c, "手机号或邮箱不能为空！")
		return
	}
	typ := req.Type
	if typ == 0 {
		// 判断账号类型
		if utils.AccountType(account) == "phone" {

			// 查询当前账号是否存在
			user, err := s.uc.FindUserByPhoneNumber(ctx, account)
			if err != nil {
				response.Error(c, err)
				return
			}

			if user != nil {
				response.Error(c, "该手机号已被注册！")
				return
			}

			// 验证验证码
			code := req.Code

			err = s.smsc.VerifyCode(&biz.Sms{
				Ctx:       ctx,
				Code:      code,
				Account:   account,
				Type:      "phone",
				SceneCode: PHONE_REGISTER_SCENE_CODE,
			})

			if err != nil {
				response.Error(c, err)
				return
			}

			// 输入密码
			pwd := req.Password
			// 确认密码
			repwd := req.RePassword

			if pwd != repwd {
				response.Error(c, "两次密码不一致，请重新输入！")
				return
			}
			// 创建用户

			// 密码加盐
			salt := s.Config.Mysql.Salt
			hashedPassword, err := hashed.Password(pwd, salt)
			if err != nil {
				response.Error(c, "系统错误，创建密码失败！")
				return
			}
			user = &biz.User{
				PhoneNumber: account,
				Salt:        salt,
				Password:    hashedPassword,
			}
			err = s.uc.Register(ctx, user)
			if err != nil {
				response.Error(c, err)
				return
			}
			response.Success(c, "注册成功！", user)
			return
		} else if utils.AccountType(account) == "email" {
			response.Success(c, "邮箱注册", nil)
			return
		}
		response.Error(c, "非法手机号或邮箱！")
		return

	} else {
		// 账号密码注册
	}
}

// 发送注册验证码
func (s *Context) SendRegisterCode(c *gin.Context) {
	ctx := c.Request.Context()

	req := biz.SmsCode{}
	// 使用 ShouldBindJSON 方法来解析JSON参数
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, err)
		return
	}

	account := req.Account
	if strings.TrimSpace(account) == "" {
		response.Error(c, "手机号或邮箱不能为空！")
		return
	}

	phoneCode := sms.GenerateRandomCode(4)

	// 判断账号类型

	clientIP := c.ClientIP()
	if utils.AccountType(account) == "phone" {
		err := s.smsc.SendSms(&biz.Sms{
			Ctx:       ctx,
			Code:      phoneCode,
			Account:   account,
			Type:      "phone",
			ClientIP:  clientIP,
			SceneCode: PHONE_REGISTER_SCENE_CODE,
		})

		if err != nil {
			response.Error(c, err)
			return
		}
		response.Success(c, "发送成功！", nil)
		return
	} else if utils.AccountType(account) == "email" {
		response.Success(c, "邮箱注册", nil)
		return
	}
	response.Error(c, "非法手机号或邮箱！")
}

// 用户登录
func (s *Context) Login(c *gin.Context) {

	ctx := c.Request.Context()

	var req biz.Login

	// 使用 ShouldBindJSON 方法来解析JSON参数
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	account := req.Account
	password := req.Password

	if strings.TrimSpace(account) == "" {
		response.Error(c, "账号不能为空！")
		return
	}

	user, err := s.uc.FindUserByAccount(ctx, account)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	if strings.TrimSpace(password) == "" {
		response.Error(c, "密码不能为空！")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password+user.Salt)); err != nil {
		response.Error(c, "密码错误！")
		return
	}

	token, err := s.uc.Token(ctx, user)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, "登录成功！", token)
}

// 校验用户信息
func (s *Context) AuthMiddleware(c *gin.Context) {
	token, err := s.uc.ValidationBearerToken(c.Request)
	if err != nil {
		response.Error(c, response.ErrAuth)
		return
	}
	c.Set("userId", token.GetUserID())
	c.Next()
}
