package utils

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 正则表达式模式
const phonePattern = "^(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\\d{8}$"
const emailPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"

// 编译正则表达式
var phoneRegexp = regexp.MustCompile(phonePattern)
var emailRegexp = regexp.MustCompile(emailPattern)

func AccountType(account string) string {

	// 判断账号类型
	if phoneRegexp.MatchString(account) {
		return "phone"
	} else if emailRegexp.MatchString(account) {
		return "email"
	}
	return "account"
}

func UserID(c *gin.Context) (int64, error) {
	userIdI, exist := c.Get("userId")
	if !exist {
		return 0, errors.New("用户登录已失效")
	}

	userId, err := strconv.ParseInt(userIdI.(string), 10, 64)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

// pageSzie为0则不显示分页
func ParsePagination(c *gin.Context) (current, pageSize int) {

	page := c.Param("page")
	size := c.Param("pageSize")

	var err error

	current, err = strconv.Atoi(page)
	if err != nil {
		current = 1
	}
	pageSize, err = strconv.Atoi(size)
	if err != nil {
		pageSize = 10
	}

	return
}
