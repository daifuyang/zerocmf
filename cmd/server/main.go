package main

import (
	"flag"
	"os"
	"zerocmf/configs"
	"zerocmf/internal/data"
	"zerocmf/internal/svc"

	"github.com/gin-gonic/gin"

	"gopkg.in/yaml.v3"
)

var configFile = flag.String("f", "configs/config.yaml", "the config file")

type App struct {
	Engine *gin.Engine
}

func (app *App) Run() error {
	app.Engine.Run()
	return nil
}

func newApp(r *gin.Engine, data *data.Data) App {
	r.Run(":8080")
	return App{
		Engine: r,
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

	svcCtx := svc.ServiceContext{
		Config: config,
	}

	// todo 初始化日志服务
	app, cleanup, err := wireApp(&svcCtx)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
