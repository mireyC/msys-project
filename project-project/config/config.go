package config

import (
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"mirey7/project-common/logs"
	"net"
	"net/http"
	"os"
	"strings"
)

var C = InitConfig()

type Config struct {
	viper *viper.Viper
	SC    *ServerConfig
	GC    *GrpcConfig
	EC    *EtcdConfig
	MC    *MysqlConfig
	JC    *JwtConfig
}

type JwtConfig struct {
	AccessExp     int
	RefreshExp    int
	AccessSecret  string
	RefreshSecret string
}

type MysqlConfig struct {
	UserName string
	Password string
	Host     string
	Port     int
	Db       string
}

type EtcdConfig struct {
	Addr     string
	UserName string
	Password string
	DialTime int64
	Weight   int64
}

type ServerConfig struct {
	Name string
	Addr string
}

type GrpcConfig struct {
	Name        string
	Addr        string
	Port        string
	Weight      int64
	UserService string
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
	conf.ReadMysqlConfig()
	conf.ReadJwtConfig()
	return conf
}

func (c *Config) ReadMysqlConfig() {
	mc := &MysqlConfig{}
	mc.UserName = c.viper.GetString("mysql.username")
	mc.Password = c.viper.GetString("mysql.password")
	mc.Host = c.viper.GetString("mysql.host")
	mc.Port = c.viper.GetInt("mysql.port")
	mc.Db = c.viper.GetString("mysql.db")
	c.MC = mc
}

func (c *Config) ReadJwtConfig() {
	jc := &JwtConfig{}
	jc.RefreshExp = c.viper.GetInt("jwt.refreshExp")
	jc.AccessExp = c.viper.GetInt("jwt.accessExp")
	jc.AccessSecret = c.viper.GetString("jwt.accessSecret")
	jc.RefreshSecret = c.viper.GetString("jwt.refreshSecret")
	c.JC = jc
}

func (c *Config) ReadServerConfig() {
	sc := &ServerConfig{}
	sc.Name = c.viper.GetString("server.name")
	sc.Addr = c.viper.GetString("server.addr")
	c.SC = sc
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

func (c *Config) ReadRedisConfig() *redis.Options {
	return &redis.Options{
		Addr:     c.viper.GetString("redis.host") + ":" + c.viper.GetString("redis.port"),
		Password: c.viper.GetString("redis.password"),
		DB:       c.viper.GetInt("redis.db"),
	}
}

// GetOutboundIP 获得对外发送消息的 IP 地址
func GetOutboundIP() string {
	// DNS 的地址，国内可以用 114.114.114.114
	conn, err := net.Dial("udp", "114.114.114.114:80")
	if err != nil {
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

// GetPublicIP 获取本地服务器的公网 IP 地址
func GetPublicIP() string {
	resp, err := http.Get("http://ipinfo.io/ip")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(ip))
}

func (c *Config) ReadGrpcConfig() {
	gc := &GrpcConfig{}
	gc.Name = c.viper.GetString("grpc.name")
	gc.Addr = GetOutboundIP() + c.viper.GetString("grpc.port")
	gc.Port = c.viper.GetString("grpc.port")
	log.Printf("grpc addr: %v", gc.Addr)
	gc.Weight = c.viper.GetInt64("grpc.weight")
	gc.UserService = c.viper.GetString("grpc.userService")
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
