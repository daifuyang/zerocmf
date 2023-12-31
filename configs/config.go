package configs

import (
	"fmt"
)

// Config 结构体用于映射配置文件中的数据
type Config struct {
	Name   string       `yaml:"Name"`
	Host   string       `yaml:"Host"`
	Port   int          `yaml:"Port"`
	Debug  bool         `yaml:"Debug"`
	Mysql  MysqlConfig  `yaml:"Mysql"`
	Redis  RedisConfig  `yaml:"Redis"`
	Sms    SmsConfig    `yaml:"Sms"`
	Oauth2 OAuth2Config `yaml:"Oauth2"`
}

// MysqlConfig 结构体用于映射MySQL配置
type MysqlConfig struct {
	Host     string `yaml:"Host"`
	Database string `yaml:"Database"`
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
	Port     int    `yaml:"Port"`
	Charset  string `yaml:"Charset"`
	Prefix   string `yaml:"Prefix"`
	Salt     string `yaml:"Salt"`
}

// DSN 格式化
func (c *MysqlConfig) Dsn(table ...bool) string {
	if len(table) > 0 && table[0] {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			c.Username,
			c.Password,
			c.Host,
			c.Port,
			c.Database,
			c.Charset)

		return dsn
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/",
		c.Username,
		c.Password,
		c.Host,
		c.Port)

	return dsn
}

// RedisConfig 结构体用于映射Redis配置
type RedisConfig struct {
	Enabled    bool   `yaml:"Enabled"`
	Addr       string `yaml:"Addr"`
	Db         int    `yaml:"Db"`
	Password   string `yaml:"Password"`
	Expiration int64  `yaml:"Expiration"`
}

// 生成结构体
type SmsConfig struct {
	AccessKeyId     string `yaml:"AccessKeyId"`
	AccessKeySecret string `yaml:"AccessKeySecret"`
	Provider        string `yaml:"Provider"`
	SignName        string `yaml:"SignName"`
	TemplateCode    string `yaml:"TemplateCode"`
	ExpiresTime     int64  `yaml:"ExpiresTime"`
}

type OAuth2Config struct {
	ClientID      string   `yaml:"ClientID"`
	ClientSecret  string   `yaml:"ClientSecret"`
	Scopes        []string `yaml:"Scopes"`
	AuthServerURL string   `yaml:"authServerURL"`
	RedirectURL   string   `yaml:"RedirectURL"`
	AuthURL       string   `yaml:"AuthURL"`
	TokenURL      string   `yaml:"TokenURL"`
}
