package main

import (
	"flag"
	"os"
	"zerocmf/configs"
	"zerocmf/internal/server"

	"gopkg.in/yaml.v3"
)

var configFile = flag.String("f", "configs/config.yaml", "the config file")

type App struct {
	*server.Server
}

func (app *App) Run(addr string) error {
	app.Server.Start() // 启动服务
	return app.Router.Run(addr)
}

func newApp(s *server.Server) App {
	return App{
		Server: s,
	}
}

// 解析配置文件
func mustLoad(configFile string, config *configs.Config) {
	// 解析配置项
	data, err := os.ReadFile(configFile)
	if err != nil {
		panic("读取配置文件失败")
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		panic("解析配置文件失败")
	}
}

func main() {
	flag.Parse()

	var config configs.Config
	mustLoad(*configFile, &config)

	// todo 初始化日志服务
	app, cleanup, err := wireApp(&config)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(":8080"); err != nil {
		panic(err)
	}
}
