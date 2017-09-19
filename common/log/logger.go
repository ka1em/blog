package log

import (
	"blog/common/setting"

	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var ZapLog *zap.SugaredLogger

func init() {
	zapCfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapLogLevel(setting.RUN_MODE)),
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
		OutputPaths:      []string{setting.LOG_OUTPUT},
		ErrorOutputPaths: []string{setting.LOG_OUTPUT},
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
	case setting.DEV_MODE:
		return zapcore.DebugLevel
	case setting.TEST_MODE:
		return zapcore.DebugLevel
	case setting.PROD_MODE:
		return zapcore.WarnLevel
	default:
		log.Fatal("error run mode")
	}
	return -1
}
