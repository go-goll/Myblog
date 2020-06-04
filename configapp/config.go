package configapp

import (
	"fmt"
	"github.com/lanfeng6/Myblog/apicontext"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Debug  bool                   `yaml:"debug"`
	Server *apicontext.ServerConf `yaml:"server"`
	Redis  *apicontext.RedisConf  `yaml:"redis"`
	MySql  *apicontext.MySqlConf  `yaml:"mysql"`
}

func (c *Config) GetMySqlConf() *apicontext.MySqlConf {
	return c.MySql
}
func (c *Config) GetServer() *apicontext.ServerConf {
	return c.Server
}
func (c *Config) GetRedisConf() *apicontext.RedisConf {
	return c.Redis
}

func (c *Config) Load(path string) {
	// 检查 path 的有效性，包含路径存在，文件 read 权限
	if path == "" {
		fmt.Println("path 不能为空")
		os.Exit(1)
	}

	// todo 路径存在，文件 read 权限

	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("解析配置文件失败", err.Error())
		os.Exit(1)
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		fmt.Println("反序列化配置文件失败", err.Error())
		os.Exit(1)
	}
}

// 定义 redis 的 key 的路径区分
const (
	REDIS_KEY_LOCAL_TOKEN_PREFIX    = "SSO:LOCAL:TOKEN"
	REDIS_KEY_LOCAL_MSGCONSUMER_KEY = "MSGCONSUMER:LOCAL:GROUPMSG"
)
