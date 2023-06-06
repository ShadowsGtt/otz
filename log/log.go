package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"sync"
)

// FormatType 日志格式类型
type FormatType string

const (
	FormatTypeConsole FormatType = "console"
	FormatTypeJSON    FormatType = "json"
)

// Config 配置
type Config struct {
	FileName   string     `yaml:"file_name"`   // 文件名
	Level      string     `yaml:"level"`       // 日志等级 debug info warn error fatal
	MaxSize    int        `yaml:"max_size"`    // 文件大小限制 MB
	MaxAge     int        `yaml:"max_age"`     // 保留天数
	MaxBackups int        `yaml:"max_backups"` // 文件数
	Compress   bool       `yaml:"compress"`    // 是否压缩
	FormatType FormatType `yaml:"format_type"` // 格式换类型
}

func newEncoder(c Config) zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "caller",
		FunctionKey:    "",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.999"),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	if c.FormatType == FormatTypeJSON {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// Levels 配置日志等级 -> zapcore.Level
var Levels = map[string]zapcore.Level{
	"":      zapcore.DebugLevel,
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
	"fatal": zapcore.FatalLevel,
}

// NewZapLog 创建zap日志
func NewZapLog(c Config) *zap.Logger {
	syncers := []zapcore.WriteSyncer{}
	// 是否写入文件
	if c.FileName != "" {
		logWriter := &lumberjack.Logger{
			Filename:   c.FileName,
			MaxSize:    c.MaxSize,    // 文件大小限制
			MaxBackups: c.MaxBackups, // 文件数
			MaxAge:     c.MaxAge,     // 保留天数
			LocalTime:  true,         // 是否本地时间 默认UTC
			Compress:   c.Compress,   // 是否压缩
		}
		syncers = append(syncers, zapcore.AddSync(logWriter))
	} else {
		// 否则写入控制台
		syncers = append(syncers, zapcore.AddSync(os.Stdout))
	}
	ws := zapcore.NewMultiWriteSyncer(syncers...)
	core := zapcore.NewCore(newEncoder(c), ws, zap.NewAtomicLevelAt(Levels[c.Level]))
	return zap.New(core,
		zap.AddCaller(), zap.AddCallerSkip(1),
	)
}

var (
	defaultConfig = Config{
		FileName:   "server.log",
		MaxSize:    100,
		MaxAge:     7,
		MaxBackups: 10,
		Compress:   false,
		FormatType: FormatTypeConsole,
	}
	defaultLogger = NewZapLog(defaultConfig)
	mutex         sync.RWMutex
)

// SetDefaultLogger 设置默认日志
func SetDefaultLogger(cfg Config) {
	mutex.Lock()
	defer mutex.Unlock()
	defaultConfig = cfg
	defaultLogger = NewZapLog(cfg)
}

// GetDefaultLogger 获取默认日志
func GetDefaultLogger() *zap.Logger {
	mutex.Lock()
	defer mutex.Unlock()
	return defaultLogger
}

func makeMsg(args ...interface{}) string {
	msg := fmt.Sprint(args...)
	return msg
}

func makeMsgFormat(format string, args ...interface{}) string {
	msg := fmt.Sprintf(format, args...)
	return msg
}

// With 设置用户自定义字段
func With(fields ...string) *zap.Logger {
	zapFields := make([]zap.Field, len(fields)/2)
	for i := range zapFields {
		zapFields[i] = zap.String(fields[2*i], fields[2*i+1])
	}
	return defaultLogger.With(zapFields...)
}

// Debug without format
func Debug(args ...interface{}) {
	defaultLogger.Debug(makeMsg(args...))
}

// Debugf without format
func Debugf(format string, args ...interface{}) {
	defaultLogger.Debug(makeMsgFormat(format, args...))
}

// Info without format
func Info(args ...interface{}) {
	defaultLogger.Info(makeMsg(args...))
}

// Infof with format
func Infof(format string, args ...interface{}) {
	defaultLogger.Info(makeMsgFormat(format, args...))
}

// Warn without format
func Warn(args ...interface{}) {
	defaultLogger.Warn(makeMsg(args...))
}

// Warnf with format
func Warnf(format string, args ...interface{}) {
	defaultLogger.Warn(makeMsgFormat(format, args...))
}

// Error without format
func Error(args ...interface{}) {
	defaultLogger.Error(makeMsg(args...))
}

// Errorf  with format
func Errorf(format string, args ...interface{}) {
	defaultLogger.Error(makeMsgFormat(format, args...))
}

// Fatal without format
func Fatal(args ...interface{}) {
	defaultLogger.Fatal(makeMsg(args...))
}

// Fatalf with format
func Fatalf(format string, args ...interface{}) {
	defaultLogger.Fatal(makeMsg(args...))
}
