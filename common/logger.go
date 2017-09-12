package common

import (
	"encoding/json"
	"io/ioutil"

	"go.uber.org/zap"
)

var Suggar *zap.SugaredLogger

func init() {
	// zap log config
	b, err := ioutil.ReadFile("/Users/ka1em/go/src/blog.ka1em.site/config/zap.json")
	if err != nil {
		panic(err)
	}

	var cfg zap.Config
	if err := json.Unmarshal(b, &cfg); err != nil {
		panic(err)
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(err.Error())
	}
	defer logger.Sync()

	Suggar = logger.Sugar()
	defer Suggar.Sync()
}
