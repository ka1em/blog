package setting

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/ini.v1"
)

var (
	// RunMode 运行模式
	RunMode string
	// SSLMode ssl模式
	SSLMode bool
	// CertFile 证书
	CertFile string
	// KeyFile 证书
	KeyFile string
	// DBHost 数据库主机
	DBHost string
	// DBPort 数据库端口
	DBPort string
	// DBUser 数据库用户
	DBUser string
	// DBPass 数据库摩玛
	DBPass string
	// DBBase 数据库
	DBBase string
	// DBParm 数据库参数
	DBParm string
	// LogPath 日志路径
	LogPath string
	// AppPath 运行路径
	AppPath string
	// TLSMinVersion min version
	TLSMinVersion string
)

const (
	// DevMode dev model
	DevMode = "dev"
	// TestMode test model
	TestMode = "test"
	// ProdMode prod model
	ProdMode = "prod"
)

// execPath 返回执行的路径
func execPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	return filepath.Abs(file)
}

func init() {
	var err error
	if AppPath, err = execPath(); err != nil {
		log.Fatal(2, "Fail to get app path: %v\n", err)
	}
}

// WorkDir 返回绝对路径
func WorkDir() (string, error) {
	wd := os.Getenv("BLOG_WORK_DIR")
	if len(wd) > 0 {
		return wd, nil
	}

	i := strings.LastIndex(AppPath, "/")
	if i == -1 {
		return AppPath, nil
	}
	return AppPath[:i], nil
}

// NewContext  init the configure
func NewContext(file string) {
	cfg, err := ini.InsensitiveLoad(file) //字段名忽略大小写
	if err != nil {
		log.Fatal(err.Error())
	}

	secServer, err := cfg.GetSection("server")
	if err != nil {
		log.Fatal(err.Error())
	}

	RunMode = secServer.Key("RUN_MODE").MustString("dev")

	SSLMode = secServer.Key("ssl_on").MustBool(false)
	CertFile = secServer.Key("cert_file").MustString("")
	KeyFile = secServer.Key("key_file").MustString("")

	secSQL, err := cfg.GetSection("mysql")
	if err != nil {
		log.Fatal(err.Error())
	}

	DBHost = secSQL.Key("DB_HOST").MustString("127.0.0.1")
	DBPort = secSQL.Key("DB_PORT").MustString("3306")
	DBUser = secSQL.Key("DB_USER").MustString("root")
	DBPass = secSQL.Key("DB_PASS").MustString("")
	DBBase = secSQL.Key("DB_BASE").MustString("lgwd")
	DBParm = secSQL.Key("DB_PARM").MustString("charset=utf8mb4&parseTime=True&loc=Local")

	secZap, err := cfg.GetSection("zap_log")
	if err != nil {
		log.Fatal(err.Error())
	}

	LogPath = secZap.Key("LOG_OUTPUT").MustString("stdout")

	TLSMinVersion = secServer.Key("TLSMinVersion").MustString("TLS12")
}
