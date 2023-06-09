package log

import (
	"context"
	"errors"
	"sync"
)

const (
	OtzLoggerKey = "OTZ_LOG_CTX_KEY"
)

var (
	defaultLogger Logger
	mutex         sync.RWMutex
)

// 默认控制台输出
func init() {
	SetDefault(Config{OutputType: OutputTypeConsole})
}

// Parser 解析器
type Parser interface {
	Parse(dst interface{}) error
}

// SetDefaultWithLoader 设置默认日志
func SetDefaultWithLoader(parser Parser) error {
	if parser == nil {
		return errors.New("loader is nil")
	}
	cfgs := []Config{}
	err := parser.Parse(&cfgs)
	if err != nil {
		return err
	}
	defaultLogger = NewZapLog(cfgs...)

	return nil
}

// SetDefault 设置默认日志
func SetDefault(cfgs ...Config) {
	mutex.Lock()
	defer mutex.Unlock()
	defaultLogger = NewZapLog(cfgs...)
}

// GetDefaultLogger 获取默认日志
func GetDefaultLogger() Logger {
	mutex.Lock()
	defer mutex.Unlock()
	return defaultLogger
}

// With 设置用户自定义字段
func With(fields ...string) Logger {
	return GetDefaultLogger().With(fields...)
}

// Context 日志记录上下文
type Context struct {
	logger  Logger
	context context.Context
}

// SetLogger 设置logger
func (ctx *Context) SetLogger(logger Logger) {
	ctx.logger = logger
}

// GetLogger 获取logger
func (ctx *Context) GetLogger() Logger {
	return ctx.logger
}

// GetCtx 获取ctx
func (ctx *Context) GetCtx() context.Context {
	return ctx.context
}

// WithCtx 通过ctx设置用户自定义字段
func WithCtx(ctx context.Context, fields ...string) context.Context {
	logCtx := GetLogCtx(ctx)
	logger := logCtx.GetLogger()
	if logger != nil {
		logger = logCtx.logger.With(fields...)
	} else {
		logger = GetDefaultLogger().With(fields...)
	}
	logCtx.SetLogger(logger)

	return ctx
}

// NewLogCtx 新建log ctx
func NewLogCtx(ctx context.Context) *Context {
	newCtx := &Context{}
	ctx = context.WithValue(ctx, OtzLoggerKey, newCtx)
	newCtx.context = ctx
	return newCtx

}

// GetLogCtx 获取ctx和logCtx
func GetLogCtx(ctx context.Context) *Context {
	// 从ctx中获取
	logCtx, ok := ctx.Value(OtzLoggerKey).(*Context)
	if ok && logCtx != nil {
		return logCtx
	}

	return NewLogCtx(ctx)
}

// Debug without format
func Debug(args ...interface{}) {
	GetDefaultLogger().Debug(makeMsg(args...))
}

// Debugf without format
func Debugf(format string, args ...interface{}) {
	GetDefaultLogger().Debug(makeMsgFormat(format, args...))
}

// DebugCtx without format
func DebugCtx(ctx context.Context, args ...interface{}) {
	logCtx, ok := ctx.Value(OtzLoggerKey).(*Context)
	if ok && logCtx != nil && logCtx.logger != nil {
		logCtx.logger.Debug(args...)
		return
	}
	GetDefaultLogger().Debug(args...)
}

// DebugCtxf with format
func DebugCtxf(ctx context.Context, format string, args ...interface{}) {
	logCtx, ok := ctx.Value(OtzLoggerKey).(*Context)
	if ok && logCtx != nil && logCtx.logger != nil {
		logCtx.logger.Debugf(format, args...)
		return
	}
	GetDefaultLogger().Debugf(format, args...)
}

// Info without format
func Info(args ...interface{}) {
	GetDefaultLogger().Info(makeMsg(args...))
}

// Infof with format
func Infof(format string, args ...interface{}) {
	GetDefaultLogger().Info(makeMsgFormat(format, args...))
}

// InfoCtx without format
func InfoCtx(ctx context.Context, args ...interface{}) {
	logCtx, ok := ctx.Value(OtzLoggerKey).(*Context)
	if ok && logCtx != nil && logCtx.logger != nil {
		logCtx.logger.Info(args...)
		return
	}
	GetDefaultLogger().Info(args...)
}

// InfoCtxf with format
func InfoCtxf(ctx context.Context, format string, args ...interface{}) {
	logCtx, ok := ctx.Value(OtzLoggerKey).(*Context)
	if ok && logCtx != nil && logCtx.logger != nil {
		logCtx.logger.Infof(format, args...)
		return
	}
	GetDefaultLogger().Infof(format, args...)
}

// Warn without format
func Warn(args ...interface{}) {
	GetDefaultLogger().Warn(makeMsg(args...))
}

// Warnf with format
func Warnf(format string, args ...interface{}) {
	GetDefaultLogger().Warn(makeMsgFormat(format, args...))
}

// WarnCtxf with format
func WarnCtxf(ctx context.Context, format string, args ...interface{}) {
	logCtx, ok := ctx.Value(OtzLoggerKey).(*Context)
	if ok && logCtx != nil && logCtx.logger != nil {
		logCtx.logger.Warnf(format, args...)
		return
	}
	GetDefaultLogger().Warnf(format, args...)
}

// WarnCtx without format
func WarnCtx(ctx context.Context, args ...interface{}) {
	logCtx, ok := ctx.Value(OtzLoggerKey).(*Context)
	if ok && logCtx != nil && logCtx.logger != nil {
		logCtx.logger.Warn(args...)
		return
	}
	GetDefaultLogger().Warn(args...)
}

// Error without format
func Error(args ...interface{}) {
	GetDefaultLogger().Error(makeMsg(args...))
}

// Errorf  with format
func Errorf(format string, args ...interface{}) {
	GetDefaultLogger().Error(makeMsgFormat(format, args...))
}

// ErrorCtx without format
func ErrorCtx(ctx context.Context, args ...interface{}) {
	logCtx, ok := ctx.Value(OtzLoggerKey).(*Context)
	if ok && logCtx != nil && logCtx.logger != nil {
		logCtx.logger.Error(args...)
		return
	}
	GetDefaultLogger().Error(args...)
}

// ErrorCtxf with format
func ErrorCtxf(ctx context.Context, format string, args ...interface{}) {
	logCtx, ok := ctx.Value(OtzLoggerKey).(*Context)
	if ok && logCtx != nil && logCtx.logger != nil {
		logCtx.logger.Errorf(format, args...)
		return
	}
	GetDefaultLogger().Errorf(format, args...)
}

// Fatal without format
func Fatal(args ...interface{}) {
	GetDefaultLogger().Fatal(makeMsg(args...))
}

// Fatalf with format
func Fatalf(format string, args ...interface{}) {
	GetDefaultLogger().Fatal(makeMsg(args...))
}

// FatalCtx with format
func FatalCtx(ctx context.Context, args ...interface{}) {
	logCtx, ok := ctx.Value(OtzLoggerKey).(*Context)
	if ok && logCtx != nil && logCtx.logger != nil {
		logCtx.logger.Fatal(args...)
		return
	}
	GetDefaultLogger().Fatal(args...)
}

// FatalCtxf with format
func FatalCtxf(ctx context.Context, format string, args ...interface{}) {
	logCtx, ok := ctx.Value(OtzLoggerKey).(*Context)
	if ok && logCtx != nil && logCtx.logger != nil {
		logCtx.logger.Fatalf(format, args...)
		return
	}
	GetDefaultLogger().Fatalf(format, args...)
}
