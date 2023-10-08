package main

import (
	"database/sql"
	"flag"
	"os"
	"zerocmf/configs"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var configFile = flag.String("f", "configs/config.yaml", "the config file")

type ServiceContext struct {
	Engine *gin.Engine
	Db     *gorm.DB
	Config configs.Config
}

func (c *ServiceContext) Start() {
	r := c.Engine
	r.Run(":8080")
}

// 初始数据库连接
func newGorm(config configs.Config) (db *gorm.DB) {
	dsn := config.Mysql.Dsn()
	sqlDb, sqlErr := sql.Open("mysql", dsn)
	if sqlErr != nil {
		panic(sqlErr)
	}
	defer sqlDb.Close()
	_, sqlErr = sqlDb.Exec("CREATE DATABASE IF NOT EXISTS " + config.Mysql.Database)
	if sqlErr != nil {
		panic(sqlErr)
	}

	db, err := gorm.Open(gormMysql.Open(config.Mysql.Dsn(true)), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return
}

func newApp(config configs.Config) *ServiceContext {
	// 初始化gin服务
	r := gin.New()
	// 初始化Db数据库
	db := newGorm(config)
	return &ServiceContext{
		Engine: r,
		Db:     db,
		Config: config,
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
	e := wireApp(config)
	e.Start()
}
