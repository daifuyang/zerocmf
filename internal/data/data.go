package data

import (
	"context"
	"database/sql"
	"encoding/json"
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
var ProviderSet = wire.NewSet(NewData, NewUserRepo, NewSmsRepo, NewDeparmentRepo, NewMenuRepo, NewRoleRepo)

type Data struct {
	db     *gorm.DB
	rdb    *redis.Client
	config *configs.Config
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

	err = biz.AutoMigrate(db, config)
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
		db:     db,
		rdb:    rdb,
		config: config,
	}

	return d, func() {
		if err := d.rdb.Close(); err != nil {

		}
	}, nil
}

// redis缓存
func (d *Data) RSet(ctx context.Context, key string, value interface{}) error {
	expiration := time.Duration(168) * time.Hour
	if d.config.Redis.Expiration > 0 {
		expiration = time.Duration(d.config.Redis.Expiration) * time.Hour
	}
	jsonData, err := json.Marshal(value)
	if err != nil {
		// 处理错误
		return err
	}
	return d.rdb.Set(ctx, key, jsonData, expiration).Err()
}

// redis读取
func (d *Data) RGet(ctx context.Context, result interface{}, key string) error {
	cachedJSON, err := d.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	// 如果数据存在于缓存中，将其解码为结构体
	err = json.Unmarshal(cachedJSON, &result)
	return err
}
