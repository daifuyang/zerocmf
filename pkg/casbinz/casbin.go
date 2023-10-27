package casbinz

import (
	"zerocmf/configs"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	_ "github.com/go-sql-driver/mysql"
)

func NewAdapter(configs *configs.Config) (e *casbin.Enforcer) {
	dsn := configs.Mysql.Dsn()
	a, _ := gormadapter.NewAdapter("mysql", dsn) // Your driver and data source.
	e, _ = casbin.NewEnforcer("static/rbac_model.conf", a)
	e.LoadPolicy()
	return e
}
