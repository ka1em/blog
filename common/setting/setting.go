package setting

import (
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"path/filepath"

	"gopkg.in/ini.v1"
)

var (
	RUN_MODE string

	SSL_ON    bool
	CERT_FILE string
	KEY_FILE  string

	DB_HOST string //mysql
	DB_PORT string
	DB_USER string
	DB_PASS string
	DB_BASE string
	DB_PARM string

	LOG_OUTPUT string

	IsWindows bool
	AppPath   string
)

const (
	DEV_MODE  = "dev"
	TEST_MODE = "test"
	PROD_MODE = "prod"
)

// execPath returns the executable path.
func execPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	return filepath.Abs(file)
}

func init() {
	IsWindows = runtime.GOOS == "windows"
	//log.New(log.CONSOLE, log.ConsoleConfig{})

	var err error
	if AppPath, err = execPath(); err != nil {
		log.Fatal(2, "Fail to get app path: %v\n", err)
	}

	// Note: we don't use path.Dir here because it does not handle case
	//	which path starts with two "/" in Windows: "//psf/Home/..."
	AppPath = strings.Replace(AppPath, "\\", "/", -1)
}

// WorkDir returns absolute path of work directory.
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

func NewContext(file string) {

	//字段名忽略大小写
	cfg, err := ini.InsensitiveLoad(file)
	if err != nil {
		log.Fatal(err.Error())
	}

	secServer, err := cfg.GetSection("server")
	if err != nil {
		log.Fatal(err.Error())
	}

	SSL_ON = secServer.Key("ssl_on").MustBool(true)
	CERT_FILE = secServer.Key("cert_file").MustString("")
	KEY_FILE = secServer.Key("key_file").MustString("")

	secSql, err := cfg.GetSection("mysql")
	if err != nil {
		log.Fatal(err.Error())
	}

	DB_HOST = secSql.Key("DB_HOST").MustString("127.0.0.1")
	DB_PORT = secSql.Key("DB_PORT").MustString("3306")
	DB_USER = secSql.Key("DB_USER").MustString("root")
	DB_PASS = secSql.Key("DB_PASS").MustString("passwd")
	DB_BASE = secSql.Key("DB_BASE").MustString("lgwd")
	DB_PARM = secSql.Key("DB_PARM").MustString("charset=utf8mb4&parseTime=True&loc=Local")

	secZap, err := cfg.GetSection("zap_log")
	if err != nil {
		log.Fatal(err.Error())
	}

	LOG_OUTPUT = secZap.Key("LOG_OUTPUT").MustString("stdout")
}
