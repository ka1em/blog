package setting

import (
	"gopkg.in/ini.v1"
)

var (
	RunMode string // RunMode 运行模式

	DBHost string // DBHost 数据库主机
	DBPort string // DBPort 数据库端口
	DBUser string // DBUser 数据库用户
	DBPass string // DBPass 数据库摩玛
	DBBase string // DBBase 数据库
	DBParm string // DBParm 数据库参数

	RedisHost string
	RedisPort string

	LogPath    string // LogPath 日志路径
	SQLLogPath string // xorm日志路径

	WxAppID     string
	WxAppSecret string
)

const (
	// DevMode dev model
	DevMode = "dev"
	// TestMode test model
	TestMode = "test"
	// ProdMode prod model
	ProdMode = "prod"
)

func init() {

}

// NewContext  init the configure
func NewContext(file string) {
	cfg, err := ini.InsensitiveLoad(file) //字段名忽略大小写
	if err != nil {
		panic(err.Error())
	}

	secServer, err := cfg.GetSection("server")
	if err != nil {
		panic(err.Error())
	}

	RunMode = secServer.Key("RunMode").MustString("dev")
	LogPath = secServer.Key("LogPath").MustString("stdout")
	SQLLogPath = secServer.Key("SQLLogPath").MustString("")
	WxAppID = secServer.Key("WxAppID").MustString("")
	WxAppSecret = secServer.Key("WxAppSecret").MustString("")

	secSQL, err := cfg.GetSection("database")
	if err != nil {
		panic(err.Error())
	}

	DBHost = secSQL.Key("DBHost").MustString("127.0.0.1")
	DBPort = secSQL.Key("DBPort").MustString("3306")
	DBUser = secSQL.Key("DBUser").MustString("root")
	DBPass = secSQL.Key("DBPass").MustString("")
	DBBase = secSQL.Key("DBBase").MustString("lgwd")
	DBParm = secSQL.Key("DBParm").MustString("charset=utf8mb4&parseTime=True&loc=Local")

	RedisHost = secSQL.Key("RedisHost").MustString("")
	RedisPort = secSQL.Key("RedisPort").MustString("")
}
