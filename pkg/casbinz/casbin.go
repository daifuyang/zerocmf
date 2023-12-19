package casbinz

import (
	"os"
	"zerocmf/configs"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	_ "github.com/go-sql-driver/mysql"
)

func NewAdapter(configs *configs.Config) (e *casbin.Enforcer) {
	dsn := configs.Mysql.Dsn(true)
	a, err := gormadapter.NewAdapter("mysql", dsn, true) // Your driver and data source.
	if err != nil {
		panic("casbin adapter" + err.Error())
	}

	currentDir, err := os.Getwd()
	if err != nil {
		panic("无法获取当前工作目录!")
	}

	e, err = casbin.NewEnforcer(currentDir+"/static/rbac_model.conf", a)
	if err != nil {
		panic(" casbin NewEnforce" + err.Error())
	}

	e.LoadPolicy()
	e.SavePolicy()

	return e
}
