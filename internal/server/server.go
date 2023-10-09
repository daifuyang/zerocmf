package server

import (
	"github.com/google/wire"
)

// 注册http服务

var ProviderSet = wire.NewSet(NewHTTPServer)
