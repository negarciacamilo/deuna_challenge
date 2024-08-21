package logger

import (
	"github.com/negarciacamilo/deuna_challenge/application/defines"
	"github.com/negarciacamilo/deuna_challenge/application/domain"
	"os"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger     *zap.Logger
	level      zapcore.Level
	DebugLevel = zapcore.DebugLevel
	InfoLevel  = zapcore.InfoLevel
	ErrorLevel = zapcore.ErrorLevel
	PanicLevel = zapcore.PanicLevel
)

const (
	logTypeKey = "type"
	SubType    = "subtype"
	errorKey   = "error"
)

func init() {
	var options []zap.Option
	options = append(options, zap.AddCaller(), zap.AddCallerSkip(1))

	Logger = zap.New(newZapCore(InfoLevel), options...)
	zap.ReplaceGlobals(Logger)
	level = InfoLevel
}

func newZapCore(lvl zapcore.Level) zapcore.Core {
	config := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    "function",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     rfc3399NanoTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// jsonEncoder := zapcore.NewJSONEncoder(config)
	console := zapcore.NewConsoleEncoder(config)

	var core zapcore.Core
	if lvl >= zapcore.ErrorLevel {
		core = zapcore.NewCore(console, zapcore.Lock(zapcore.AddSync(os.Stderr)), lvl)
	} else {
		core = zapcore.NewCore(console, zapcore.Lock(zapcore.AddSync(os.Stdout)), lvl)
	}

	return core
}

func rfc3399NanoTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	const RFC3339Micro = "2006-01-02T15:04:05.000000Z07:00"

	enc.AppendString(t.UTC().Format(RFC3339Micro))
}

func Debug(message, logType string, ctx *domain.ContextInformation, t ...map[string]any) {
	if level <= DebugLevel {
		tags := getTags(t...)
		tags[logTypeKey] = logType
		Logger.Debug(message, getFields(ctx, tags)...)
	}
}

func Info(message, logType string, ctx *domain.ContextInformation, t ...map[string]any) {
	if level <= InfoLevel {
		tags := getTags(t...)
		tags[logTypeKey] = logType
		Logger.Info(message, getFields(ctx, tags)...)
	}
}

func Error(message, logType string, err error, ctx *domain.ContextInformation, t ...map[string]any) {
	if level <= ErrorLevel {
		tags := getTags(t...)
		tags[logTypeKey] = logType
		if err != nil {
			tags[errorKey] = err.Error()
		}
		Logger.Error(message, getFields(ctx, tags)...)
	}
}

func Panic(message, logType string, err error, ctx *domain.ContextInformation, t ...map[string]any) {
	tags := getTags(t...)
	tags[logTypeKey] = logType
	if err != nil {
		tags[errorKey] = err.Error()
	}
	Logger.Panic(message, getFields(ctx, tags)...)
}

func getFields(ctx *domain.ContextInformation, tags map[string]any) []zapcore.Field {
	var fields []zapcore.Field
	for k, v := range tags {
		fields = append(fields, zap.Any(k, v))
	}

	if ctx != nil {
		if ctx.RequestInfo != nil {
			if ctx.RequestInfo.RequestID != "" {
				fields = append(fields, zap.Field{Key: defines.XRequestID, Type: zapcore.StringType, String: ctx.RequestInfo.RequestID})
			}
		}
	}
	return fields
}

func getTags(tags ...map[string]any) map[string]any {
	if tags == nil {
		return map[string]any{}
	}

	if tags[0] == nil {
		tags[0] = make(map[string]any)
	}
	return tags[0]
}

func GetCallerFunctionName() string {
	pc, _, _, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()
	lastSlash := strings.LastIndexByte(funcName, '/')
	if lastSlash < 0 {
		lastSlash = 0
	}
	lastDot := strings.LastIndexByte(funcName[lastSlash:], '.') + lastSlash

	return funcName[lastDot+1:]
}
