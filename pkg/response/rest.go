package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data,omitempty"`
	Msg  string      `json:"msg"`
}

func Success(c *gin.Context, msg string, data interface{}) {
	res := response{
		Code: 1,
		Data: data,
		Msg:  msg,
	}
	c.JSON(http.StatusOK, res)
}

func Error(c *gin.Context, msg interface{}) {

	m := "服务器异常"

	res := response{
		Code: 0,
		Msg:  m,
	}

	// msg 断言

	code := http.StatusOK

	switch msg.(type) {
	case string:
		res.Msg = msg.(string)
	case error:
		err := msg.(error)
		// 身份失效
		if errors.Is(err, ErrAuth) {
			code = http.StatusUnauthorized
			res.Msg = ErrAuth.Error()
		}

		// 异常处理
		if errors.Is(err, gorm.ErrRecordNotFound) {
			res.Msg = gorm.ErrRecordNotFound.Error()
		}
		// 绑定失败则是非法请求，直接重定向到404
		if errors.Is(err, ErrBind) {
			code = http.StatusNotFound
		}
	}

	c.JSON(code, res)
	c.Abort()

}
