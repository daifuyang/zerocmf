package service

import (
	"fmt"
	"regexp"
	"strings"
	"zerocmf/internal/biz"
	"zerocmf/pkg/code"
	"zerocmf/pkg/response"

	"github.com/gin-gonic/gin"
)

// 用户注册
func (s *Context) Register(c *gin.Context) {

	req := new(biz.Register)
	c.ShouldBind(&req)

	account := req.Account
	if strings.TrimSpace(account) == "" {
		response.Error(c, "手机号或邮箱不能为空！", nil)
		return
	}
	typ := req.Type
	if typ == 0 {

		// 正则表达式模式
		const phonePattern = "^(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\\d{8}$"
		const emailPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"

		// 编译正则表达式
		phoneRegexp := regexp.MustCompile(phonePattern)
		emailRegexp := regexp.MustCompile(emailPattern)
		// 判断账号类型
		if phoneRegexp.MatchString(account) {

			// 验证验证码
			code := req.Code

			// 输入密码

			// 确认密码

			response.Success(c, "手机号注册", nil)
			return
		} else if emailRegexp.MatchString(account) {
			response.Success(c, "邮箱注册", nil)
			return
		}

		response.Error(c, "非法手机号或邮箱！", nil)
		return

	} else {
		// 账号密码注册
	}
}

// 发送注册验证码
func (s *Context) SendRegisterCode(c *gin.Context) {
	ctx := c.Request.Context()

	req := biz.Register{}
	c.ShouldBind(&req)

	account := req.Account
	if strings.TrimSpace(account) == "" {
		response.Error(c, "手机号或邮箱不能为空！", nil)
		return
	}
	// 正则表达式模式
	const phonePattern = "^(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\\d{8}$"
	const emailPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"

	// 编译正则表达式
	phoneRegexp := regexp.MustCompile(phonePattern)
	emailRegexp := regexp.MustCompile(emailPattern)

	phoneCode := code.GenerateRandomCode(4)

	// 判断账号类型
	if phoneRegexp.MatchString(account) {

		//
		err := s.smsc.SendSms(&biz.Sms{
			Ctx:  ctx,
			Code: phoneCode,
		})

		if err != nil {
			return
		}

		templateParam := fmt.Sprintf("{\"code\":\"%s\"}", phoneCode)
		// 发送短信验证码
		err = s.sms.SendSms(account, templateParam)
		if err != nil {
			response.Error(c, "发送失败！", err.Error())
			return
		}

		response.Success(c, "验证为："+templateParam, nil)
		return
	} else if emailRegexp.MatchString(account) {
		response.Success(c, "邮箱注册", nil)
		return
	}

	response.Error(c, "非法手机号或邮箱！", nil)
	return
}
