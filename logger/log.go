package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"time"
)

type Logger struct {
	filename string
	debug    bool
	*zap.Logger
}

func NewLogger(filename string, debug bool) *Logger {
	return &Logger{
		filename: filename,
		debug:    debug,
		Logger:   getLog(filename, debug),
	}
}

func getLog(filename string, debug bool) *zap.Logger {

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel
	})
	waringLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.ErrorLevel
	})
	var (
		infoWriter, waringWriter, errorWriter io.Writer
	)
	if debug {
		infoWriter, waringWriter, errorWriter = os.Stdout, os.Stdout, os.Stderr
	} else {
		infoWriter, waringWriter, errorWriter =
			getWriter(filename+"_info.log"),
			getWriter(filename+"_waring.log"),
			getWriter(filename+"_error.log")
	}

	cores := []zapcore.Core{
		getCore(infoLevel, infoWriter),
		getCore(waringLevel, waringWriter),
		getCore(errorLevel, errorWriter),
	}

	core := zapcore.NewTee(cores...)
	caller := zap.AddCaller()
	skip := zap.AddCallerSkip(0)
	//development := zap.Development()
	logger := zap.New(core, caller /*development*/, skip)
	return logger
}

func getCore(levelEnableFn zap.LevelEnablerFunc, writer io.Writer) zapcore.Core {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(writer)),
		levelEnableFn)
	return core
}

func getWriter(filename string) *lumberjack.Logger {
	today := time.Now().Format("20060102")
	return &lumberjack.Logger{
		Filename:   fmt.Sprintf("./logs/%s/%s", today, filename),
		MaxSize:    128,
		MaxBackups: 30,
		MaxAge:     7,
		Compress:   true,
	}
}
