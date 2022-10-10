package logger

import (
	"context"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/gemalto/flume"
	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	uberZap "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	once                 sync.Once
	sugaredLoggerPointer unsafe.Pointer
)

func Init(level string, fields ...interface{}) {
	once.Do(func() {
		atomic.StorePointer(&sugaredLoggerPointer, unsafe.Pointer(newZapLogger(level).With(fields...)))
	})
}

func newZapLogger(level string) *zap.SugaredLogger {
	zapLevel := getZapLevel(level)
	writer := zapcore.AddSync(colorable.NewColorableStdout())

	encCfg := uberZap.NewProductionEncoderConfig()
	encCfg.TimeKey = "time"
	encCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	flumeEncoder := zapcore.Encoder(flume.NewColorizedConsoleEncoder((*flume.EncoderConfig)(&encCfg), nil))

	core := zapcore.NewCore(flumeEncoder, writer, zapLevel)

	return zap.New(core, zap.AddCallerSkip(2), zap.AddCaller()).Sugar()
}

func L() *zap.SugaredLogger {
	if _l := (*zap.SugaredLogger)(atomic.LoadPointer(&sugaredLoggerPointer)); _l != nil {
		return _l
	}
	panic("sugared logger was not initialized")
}

func getZapLevel(level string) zapcore.Level {
	switch level {
	case zapcore.InfoLevel.String():
		return zapcore.InfoLevel
	case zapcore.WarnLevel.String():
		return zapcore.WarnLevel
	case zapcore.DebugLevel.String():
		return zapcore.DebugLevel
	case zapcore.ErrorLevel.String():
		return zapcore.ErrorLevel
	case zapcore.FatalLevel.String():
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

type fieldsKey struct{}

var _fieldsKey fieldsKey

func AddFields(ctx context.Context, newFields ...interface{}) context.Context {
	if len(newFields) == 0 {
		return ctx
	}

	value := ctx.Value(_fieldsKey)

	fields, ok := value.([]interface{})
	if ok {
		return context.WithValue(ctx, _fieldsKey, append(fields, newFields...))
	}

	return context.WithValue(ctx, _fieldsKey, newFields)
}

func GetFields(ctx context.Context) []interface{} {
	value := ctx.Value(_fieldsKey)

	fields, ok := value.([]interface{})
	if !ok {
		return nil
	}

	return fields
}
