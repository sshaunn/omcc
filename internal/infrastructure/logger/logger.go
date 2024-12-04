package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
	"time"
)

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	With(fields ...Field) Logger
	Sync() error
}

type Field = zapcore.Field

// String 提供便捷的字段创建方法
func String(key string, value string) Field {
	return zap.String(key, value)
}

func Int64(key string, value int64) Field {
	return zap.Int64(key, value)
}

func Error(err error) Field {
	return zap.Error(err)
}

func Any(key string, value interface{}) Field {
	return zap.Any(key, value)
}

func Duration(key string, value time.Duration) Field {
	return zap.Duration(key, value)
}

type zapLogger struct {
	*zap.Logger
}

func NewLogger() Logger {
	//config := zap.NewProductionConfig()
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    customLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   customCallerEncoder,
	}

	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.Encoding = "console" // 使用控制台编码器而不是JSON

	logger, err := config.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)
	if err != nil {
		panic(err)
	}

	return &zapLogger{logger}
}

// 自定义时间编码器
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// 自定义日志级别编码器
func customLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf("%-5s", l.CapitalString()))
}

// 自定义调用者编码器
func customCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	// 获取相对路径
	path := caller.TrimmedPath()
	// 提取文件名和行号
	parts := strings.Split(path, "/")
	if len(parts) > 2 {
		path = fmt.Sprintf("%s/%s", parts[len(parts)-2], parts[len(parts)-1])
	}
	enc.AppendString(fmt.Sprintf("%-20s", path))
}

// 添加新的字段格式化函数
func formatField(key string, value interface{}) Field {
	return zap.Any(fmt.Sprintf("%-15s", key), value)
}

// Debug 实现 Logger 接口
func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.Logger.Debug(msg, fields...)
}

func (l *zapLogger) Info(msg string, fields ...Field) {
	l.Logger.Info(msg, fields...)
}

func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.Logger.Warn(msg, fields...)
}

func (l *zapLogger) Error(msg string, fields ...Field) {
	l.Logger.Error(msg, fields...)
}

func (l *zapLogger) Fatal(msg string, fields ...Field) {
	l.Logger.Fatal(msg, fields...)
}

func (l *zapLogger) With(fields ...Field) Logger {
	return &zapLogger{l.Logger.With(fields...)}
}
