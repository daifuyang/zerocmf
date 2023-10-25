package response

import (
	"errors"

	"gorm.io/gorm"
)

var ErrAuth = errors.New("用户身份已失效！")

var ErrRecordNotFound = errors.New("该数据不存在或已被删除！")

var ErrBind = errors.New("数据绑定失败！") // 重定向到404

func init() {
	gorm.ErrRecordNotFound = ErrRecordNotFound
}
