package initialize

import (
	"os"

	"github.com/wangn-tech/tiny-douyin/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LoggerSetup 创建一个 zap.Logger，支持 json/console 输出以及文件/控制台输出
func LoggerSetup(c config.LogConfig) *zap.Logger {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var encoder zapcore.Encoder
	if c.Format == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	var ws zapcore.WriteSyncer
	switch c.Output {
	case "stderr":
		ws = zapcore.AddSync(os.Stderr)
	case "file":
		f, err := os.OpenFile(c.FilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			// 回退到 stdout
			ws = zapcore.AddSync(os.Stdout)
		} else {
			ws = zapcore.AddSync(f)
		}
	default:
		ws = zapcore.AddSync(os.Stdout)
	}

	level := parseLevel(c.Level)
	core := zapcore.NewCore(encoder, ws, level)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(0))
	return logger
}

func parseLevel(s string) zapcore.Level {
	switch s {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
