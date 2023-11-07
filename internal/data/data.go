package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
	"zerocmf/configs"
	"zerocmf/internal/biz"
	cmfOauth2 "zerocmf/pkg/oauth2"

	"github.com/go-oauth2/oauth2/v4/server"
	"golang.org/x/oauth2"

	"github.com/go-oauth2/mysql/v4"
	"github.com/go-redis/redis/extra/redisotel"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewUserRepo, NewSmsRepo, NewDeparmentRepo, NewMenuRepo, NewRoleRepo)

type Data struct {
	db         *gorm.DB
	rdb        *redis.Client
	srv        *server.Server
	config     *configs.Config
	oauth2Conf oauth2.Config
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

	db, err := gorm.Open(gormMysql.Open(config.Mysql.Dsn(true)), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// use mysql token store
	tokenStore := mysql.NewDefaultStore(
		mysql.NewConfig(config.Mysql.Dsn(true)),
	)

	oauth2Conf := oauth2.Config{
		ClientID:     config.Oauth2.ClientID,
		ClientSecret: config.Oauth2.ClientSecret,
		Scopes:       config.Oauth2.Scopes,
		RedirectURL:  config.Oauth2.RedirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.Oauth2.AuthURL,
			TokenURL: config.Oauth2.TokenURL,
		},
	}

	srv := cmfOauth2.NewServer(config, tokenStore)

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
		db:         db,
		rdb:        rdb,
		config:     config,
		srv:        srv,
		oauth2Conf: oauth2Conf,
	}

	return d, func() {
		tokenStore.Close()
		d.rdb.Close()
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
