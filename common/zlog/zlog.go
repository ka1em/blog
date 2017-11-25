package zlog

import (
	"blog/common/setting"
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapLog log
var ZapLog *zap.SugaredLogger

// ZapLogInit 初始化
func ZapLogInit() {
	if setting.LogPath != "stdout" {
		fd, err := os.OpenFile(setting.LogPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err.Error())
		}
		defer fd.Close()
	}

	zapCfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapLogLevel(setting.RunMode)),
		Development: true,
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "T",
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "C",
			MessageKey:     "M",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{setting.LogPath},
		ErrorOutputPaths: []string{setting.LogPath},
	}

	zaplog, err := zapCfg.Build()
	if err != nil {
		panic(err.Error())
	}
	defer zaplog.Sync()

	ZapLog = zaplog.Sugar()
	defer ZapLog.Sync()
}

func zapLogLevel(runMode string) zapcore.Level {
	switch runMode {
	case setting.DevMode:
		return zapcore.DebugLevel
	case setting.TestMode:
		return zapcore.DebugLevel
	case setting.ProdMode:
		return zapcore.InfoLevel
	default:
		log.Fatal("error run mode")
		return 0
	}
}
