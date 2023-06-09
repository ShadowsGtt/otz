package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

// FormatType 日志格式类型
type FormatType string

const (
	FormatTypeText FormatType = "text" // 文本格式
	FormatTypeJSON FormatType = "json" // json格式
)

// OutputType 日志输出类型
type OutputType string

const (
	OutputTypeConsole OutputType = "console" // 控制台
	OutputTypeFile    OutputType = "file"    // 文件
)

// Config 配置
type Config struct {
	OutputType OutputType `yaml:"output_type"` // 输出位置 console/file
	FileName   string     `yaml:"file_name"`   // 文件名
	Level      string     `yaml:"level"`       // 日志等级 debug info warn error fatal
	MaxSize    int        `yaml:"max_size"`    // 文件大小限制 MB
	MaxAge     int        `yaml:"max_age"`     // 保留天数
	MaxBackups int        `yaml:"max_backups"` // 文件数
	Compress   bool       `yaml:"compress"`    // 是否压缩
	FormatType FormatType `yaml:"format_type"` // 格式换类型 text/json
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

// ZapLog zap日志
type ZapLog struct {
	zapLog *zap.Logger
	cfgs   []Config
}

// NewZapLog 创建zap日志
func NewZapLog(cfgs ...Config) Logger {
	cores := []zapcore.Core{}
	for _, c := range cfgs {
		cores = append(cores, createZapCore(c))
	}
	return &ZapLog{
		zapLog: zap.New(
			zapcore.NewTee(cores...),
			zap.AddCaller(),
			zap.AddCallerSkip(1),
		),
		cfgs: cfgs,
	}
}

func createZapCore(c Config) zapcore.Core {
	wr := getOutputWriter(c)
	ws := zapcore.AddSync(wr)
	encoder := newEncoder(c)
	level := zap.NewAtomicLevelAt(Levels[c.Level])
	core := zapcore.NewCore(encoder, ws, level)

	return core
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

func getOutputWriter(c Config) io.Writer {
	switch c.OutputType {
	case OutputTypeFile:
		return &lumberjack.Logger{
			Filename:   c.FileName,
			MaxSize:    c.MaxSize,    // 文件大小限制
			MaxBackups: c.MaxBackups, // 文件数
			MaxAge:     c.MaxAge,     // 保留天数
			LocalTime:  true,         // 是否本地时间 默认UTC
			Compress:   c.Compress,   // 是否压缩
		}
	default:
		return zapcore.Lock(os.Stdout)
	}
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
func (zl *ZapLog) With(fields ...string) Logger {
	zapFields := make([]zap.Field, len(fields)/2)
	for i := range zapFields {
		zapFields[i] = zap.Any(fields[2*i], fields[2*i+1])
	}
	return &ZapLog{
		zapLog: zl.zapLog.With(zapFields...),
		cfgs:   zl.cfgs,
	}
}

// Debug without format
func (zl *ZapLog) Debug(args ...interface{}) {
	zl.zapLog.Debug(makeMsg(args...))
}

// Debugf without format
func (zl *ZapLog) Debugf(format string, args ...interface{}) {
	zl.zapLog.Debug(makeMsgFormat(format, args...))
}

// Info without format
func (zl *ZapLog) Info(args ...interface{}) {
	zl.zapLog.Info(makeMsg(args...))
}

// Infof with format
func (zl *ZapLog) Infof(format string, args ...interface{}) {
	zl.zapLog.Info(makeMsgFormat(format, args...))
}

// Warn without format
func (zl *ZapLog) Warn(args ...interface{}) {
	zl.zapLog.Warn(makeMsg(args...))
}

// Warnf with format
func (zl *ZapLog) Warnf(format string, args ...interface{}) {
	zl.zapLog.Warn(makeMsgFormat(format, args...))
}

// Error without format
func (zl *ZapLog) Error(args ...interface{}) {
	zl.zapLog.Error(makeMsg(args...))
}

// Errorf  with format
func (zl *ZapLog) Errorf(format string, args ...interface{}) {
	zl.zapLog.Error(makeMsgFormat(format, args...))
}

// Fatal without format
func (zl *ZapLog) Fatal(args ...interface{}) {
	zl.zapLog.Fatal(makeMsg(args...))
}

// Fatalf with format
func (zl *ZapLog) Fatalf(format string, args ...interface{}) {
	zl.zapLog.Fatal(makeMsg(args...))
}

// Flush with format
func (zl *ZapLog) Flush() error {
	return zl.zapLog.Sync()
}
