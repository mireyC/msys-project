package config

import (
	"github.com/spf13/viper"
	"log"
	"mirey7/project-common/logs"
	"os"
)

var C = InitConfig()

type Config struct {
	viper *viper.Viper
	SC    *ServerConfig
	EC    *EtcdConfig
	GC    *GrpcConfig
}

type ServerConfig struct {
	Name string
	Addr string
}

type GrpcConfig struct {
	UserName string
}

type EtcdConfig struct {
	Addr     string
	UserName string
	Password string
	DialTime int64
}

func InitConfig() *Config {
	conf := &Config{
		viper: viper.New(),
	}

	workDir, _ := os.Getwd()
	conf.viper.SetConfigName("config")
	conf.viper.SetConfigType("yaml")
	conf.viper.AddConfigPath(workDir + "/config")
	err := conf.viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}
	conf.ReadServerConfig()
	conf.InitZapLog()
	conf.ReadGrpcConfig()
	conf.ReadEtcdConfig()
	return conf
}

func (c *Config) ReadServerConfig() {
	sc := &ServerConfig{}
	sc.Name = c.viper.GetString("server.name")
	sc.Addr = c.viper.GetString("server.addr")
	c.SC = sc
}

func (c *Config) ReadGrpcConfig() {
	gc := &GrpcConfig{}
	gc.UserName = c.viper.GetString("grpc.userName")
	c.GC = gc
}

func (c *Config) ReadEtcdConfig() {
	ec := &EtcdConfig{}
	ec.Addr = c.viper.GetString("etcd.addr")
	ec.UserName = c.viper.GetString("etcd.userName")
	ec.Password = c.viper.GetString("etcd.password")
	ec.DialTime = c.viper.GetInt64("etcd.dialTime")
	c.EC = ec
}

func (c *Config) InitZapLog() {
	//从配置中读取日志配置，初始化日志
	lc := &logs.LogConfig{
		DebugFileName: c.viper.GetString("zap.debugFileName"),
		InfoFileName:  c.viper.GetString("zap.infoFileName"),
		WarnFileName:  c.viper.GetString("zap.warnFileName"),
		MaxSize:       c.viper.GetInt("maxSize"),
		MaxAge:        c.viper.GetInt("maxAge"),
		MaxBackups:    c.viper.GetInt("maxBackups"),
	}
	err := logs.InitLogger(lc)
	if err != nil {
		log.Fatalln(err)
	}
}
