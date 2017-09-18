package config

import "github.com/go-ini/ini"

const (
	CONF_DEV_PATH  = "conf/dev.cfg"
	CONF_TEST_PATH = "conf/test.cfg"
	CONF_PROD_PATH = "conf/prod.cfg"
)

func init() {
	ini.Load("./config/config")
}
