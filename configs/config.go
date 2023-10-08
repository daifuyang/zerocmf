package configs

import "fmt"

// Config 结构体用于映射配置文件中的数据
type Config struct {
	Name  string      `yaml:"Name"`
	Host  string      `yaml:"Host"`
	Port  int         `yaml:"Port"`
	Mysql MysqlConfig `yaml:"Mysql"`
	Redis RedisConfig `yaml:"Redis"`
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
	AuthCode string `yaml:"AuthCode"`
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
	Enabled  bool   `yaml:"Enabled"`
	Host     string `yaml:"Host"`
	Database int    `yaml:"Database"`
	Port     int    `yaml:"Port"`
	Password string `yaml:"Password"`
}
