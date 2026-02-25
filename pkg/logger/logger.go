package logger

import (
	"context"
	"os"
	"sync"

	"github.com/dev2choiz/api-skeleton/pkg/contextapp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	baseLogger *zap.Logger
	once       sync.Once
)

func InitLogger(logFile string, level zapcore.Level, isDev bool) {
	once.Do(func() {
		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    1, // MB
			MaxBackups: 7,
			MaxAge:     30, // days
			Compress:   true,
		})

		consoleWriter := zapcore.AddSync(os.Stdout)

		var consoleEncoder zapcore.Encoder

		if isDev {
			devCfg := zap.NewDevelopmentEncoderConfig()
			devCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
			devCfg.EncodeTime = zapcore.ISO8601TimeEncoder
			consoleEncoder = zapcore.NewConsoleEncoder(devCfg)
		} else {
			prodCfg := zap.NewProductionEncoderConfig()
			prodCfg.EncodeTime = zapcore.ISO8601TimeEncoder
			consoleEncoder = zapcore.NewJSONEncoder(prodCfg)
		}

		fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

		core := zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, consoleWriter, level),
			zapcore.NewCore(fileEncoder, fileWriter, level),
		)

		baseLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	})
}

func Get(ctx context.Context) *zap.Logger {
	if baseLogger == nil {
		panic("logger not initialized")
	}

	l := baseLogger

	cid := contextapp.GetCorrelationID(ctx)
	if cid != "" {
		l = l.With(zap.String("correlation_id", cid))
	}

	userID := contextapp.GetUser(ctx).ID
	if userID != "" {
		l = l.With(zap.String("user_id", userID))
	}

	return l
}

func GetZapLogLevel(levelStr string) zapcore.Level {
	var level zapcore.Level
	err := level.UnmarshalText([]byte(levelStr))
	if err != nil {
		return zap.InfoLevel
	}

	return level
}
