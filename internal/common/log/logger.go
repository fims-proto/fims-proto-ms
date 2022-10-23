package log

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func InitLogger() {
	config := zap.NewProductionConfig()
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig = encoderConfig

	if viper.GetBool("logger.debug") {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	if !viper.GetBool("logger.jsonEncoding") {
		config.Encoding = "console"
	}

	var err error
	logger, err = config.Build(zap.AddCallerSkip(2))
	if err != nil {
		panic(err)
	}
}

func SyncLogger() {
	_ = logger.Sync()
}

func InfoWithoutCxt(template string, fmtArgs ...any) {
	info(context.Background(), template, fmtArgs...)
}

func DebugWithoutCxt(template string, fmtArgs ...any) {
	debug(context.Background(), template, fmtArgs...)
}

func ErrWithoutCxt(err error, template string, fmtArgs ...any) {
	errLog(context.Background(), err, template, fmtArgs...)
}

func Info(ctx context.Context, template string, fmtArgs ...any) {
	info(ctx, template, fmtArgs...)
}

func Debug(ctx context.Context, template string, fmtArgs ...any) {
	debug(ctx, template, fmtArgs...)
}

func Err(ctx context.Context, err error, template string, fmtArgs ...any) {
	errLog(ctx, err, template, fmtArgs...)
}

func info(ctx context.Context, template string, fmtArgs ...any) {
	if logger != nil && logger.Core().Enabled(zap.InfoLevel) {
		logger.Info(
			getMessage(template, fmtArgs),
		)
	}
}

func debug(ctx context.Context, template string, fmtArgs ...any) {
	if logger != nil && logger.Core().Enabled(zap.DebugLevel) {
		logger.Debug(
			getMessage(template, fmtArgs),
		)
	}
}

func errLog(ctx context.Context, err error, template string, fmtArgs ...any) {
	if logger != nil && logger.Core().Enabled(zap.ErrorLevel) {
		logger.Error(
			getMessage(template, fmtArgs),
			zap.Error(err),
		)
	}
}

// getMessage format with Sprint, Sprintf, or neither.
func getMessage(template string, fmtArgs []any) string {
	if len(fmtArgs) == 0 {
		return template
	}

	if template != "" {
		return fmt.Sprintf(template, fmtArgs...)
	}

	if len(fmtArgs) == 1 {
		if str, ok := fmtArgs[0].(string); ok {
			return str
		}
	}
	return fmt.Sprint(fmtArgs...)
}
