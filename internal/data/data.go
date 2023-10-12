package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"zerocmf/configs"
	"zerocmf/internal/biz"

	"github.com/go-redis/redis/extra/redisotel"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"

	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewUserRepo, NewSmsRepo)

type Data struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewData(config *configs.Config) (*Data, func(), error) {
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

	db, err := gorm.Open(mysql.Open(config.Mysql.Dsn(true)), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = biz.AutoMigrate(db)
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password,
		DB:       int(config.Redis.Db),
	})

	// 使用 context 设置超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 执行 PING 命令
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	// 输出 PING 响应
	fmt.Println("redis ping response：" + pong)

	rdb.AddHook(redisotel.TracingHook{})

	d := &Data{
		db:  db,
		rdb: rdb,
	}

	return d, func() {
		if err := d.rdb.Close(); err != nil {

		}
	}, nil
}
