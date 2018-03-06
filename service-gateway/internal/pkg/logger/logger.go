package logger

import (
	"github.com/1ambda/go-ref/service-gateway/internal/pkg/config"
	"go.uber.org/zap"
)

var debug = true

func init() {
	s := config.GetSpecification()

	if s.Env != "LOCAL" {
		debug = false
	}
}

func Info(arg ...interface{}) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	sugar.Info(arg...)
}

func Infof(msg string, keysAndValues ...interface{}) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	sugar.Infof(msg, keysAndValues)
}

func Infow(msg string, keysAndValues ...interface{}) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	sugar.Infow(msg, keysAndValues...)
}

func Error(arg ...interface{}) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	sugar.Error(arg...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	sugar.Errorw(msg, keysAndValues...)
}

func Fatal(arg ...interface{}) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	sugar.Fatal(arg...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	sugar.Fatalw(msg, keysAndValues...)
}
