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

	DBHost string //mysql
	DBPort string
	DBUser string
	DBPass string
	DBBase string
	DBParm string

	IsWindows bool
	AppPath   string
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

	DBHost = secServer.Key("DBHost").MustString("127.0.0.1")
	DBPort = secServer.Key("DBPort").MustString("3306")
	DBUser = secServer.Key("DBUser").MustString("root")
	DBPass = secServer.Key("DBPass").MustString("passwd")
	DBBase = secServer.Key("DBBase").MustString("lgwd")
	DBParm = secServer.Key("DBParm").MustString("charset=utf8mb4&parseTime=True&loc=Local")

}
