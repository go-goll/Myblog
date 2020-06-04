package apicontext

import (
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"time"
)

type ServerConf struct {
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	Schema string `yaml:"schema"`
	Domain string `yaml:"domain"`
}
type MySqlConf struct {
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	User   string `yaml:"user"`
	Passwd string `yaml:"passwd"`
	DbName string `yaml:"dbname"`
}
type RedisConf struct {
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	Passwd string `yaml:"passwd"`
	Db     int    `yaml:"db"`
}
type RedisOption struct {
	Host   string
	Port   string
	Passwd string
	DbNum  int
}
type Context struct {
	Ds   *gorm.DB
	R    *redis.Client
	Conf IConfig
}
type IConfig interface {
	Load(string)
	GetMySqlConf() *MySqlConf
	GetRedisConf() *RedisConf
	GetServer() *ServerConf
}

func (s *ServerConf) GetServerAddr() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}

func (o *RedisOption) Addr() string {
	return fmt.Sprintf("%s:%s", o.Host, o.Port)
}

func (o *RedisOption) String() string {
	return fmt.Sprintf("[%s:%s][%s],db[%d]", o.Host, o.Port, o.Passwd, o.DbNum)
}
func NewRedisClient(opt *RedisOption) (*redis.Client, error) {
	option := &redis.Options{
		Addr:     opt.Addr(),
		Password: opt.Passwd,
		DB:       opt.DbNum,
	}

	c := redis.NewClient(option)
	_, err := c.Ping().Result()
	if err != nil {
		fmt.Println("", "ping redis[%s]失败, %s", opt.String(), err.Error())
		return nil, err
	}

	fmt.Println("", "连接 redis[%s]成功", opt.String())
	return c, nil
}

func NewContext(quit <-chan struct{}, conf IConfig) (*Context, error) {
	redisConf := conf.GetRedisConf()
	rOpt := &RedisOption{
		Host:   redisConf.Host,
		Port:   redisConf.Port,
		Passwd: redisConf.Passwd,
		DbNum:  redisConf.Db,
	}

	redisC, err := NewRedisClient(rOpt)
	if err != nil {
		fmt.Println("", "初始化认证信息的 redis 失败, %s", err.Error())
		return nil, err
	}
	dsconf := conf.GetMySqlConf()
	var ds *gorm.DB
	ds, err = gorm.Open("mysql", dsconf.User+":"+dsconf.Passwd+"@tcp("+dsconf.Host+":"+dsconf.Port+")/"+dsconf.DbName+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	ds.DB().SetMaxIdleConns(5)
	ds.DB().SetMaxOpenConns(50)
	ds.DB().SetConnMaxLifetime(time.Second * 60)
	go func() {
		select {
		case <-quit:
			fmt.Println("", "关闭业务数据库 mysql")
			ds.Close()
			redisC.Close()
		}
	}()

	ctx := &Context{
		Ds:   ds,
		R:    redisC,
		Conf: conf,
	}

	return ctx, nil
}
