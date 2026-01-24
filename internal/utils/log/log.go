package log

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger
var atomicLevel = zap.NewAtomicLevelAt(zap.InfoLevel)

var consoleEncoder = zapcore.EncoderConfig{
	TimeKey:       "time",
	LevelKey:      "level",
	MessageKey:    "msg",
	CallerKey:     "caller",
	StacktraceKey: "stacktrace",
	EncodeLevel:   zapcore.CapitalLevelEncoder,
	EncodeTime:    zapcore.RFC3339TimeEncoder,
	EncodeCaller:  zapcore.ShortCallerEncoder,
}

func init() {
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(consoleEncoder),
		zapcore.AddSync(os.Stdout),
		atomicLevel,
	)
	opts := []zap.Option{
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	}
	Logger = zap.New(core, opts...).Sugar()
}

func parseLevel(level string) zapcore.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return zapcore.DebugLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func SetLevel(level string) {
	atomicLevel.SetLevel(parseLevel(level))
}

func Infof(template string, args ...interface{}) {
	Logger.Infof(template, args...)
}

func Errorf(template string, args ...interface{}) {
	Logger.Errorf(template, args...)
}

func Warnf(template string, args ...interface{}) {
	Logger.Warnf(template, args...)
}

func Debugf(template string, args ...interface{}) {
	Logger.Debugf(template, args...)
}
