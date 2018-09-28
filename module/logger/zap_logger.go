package logger

import (
	"github.com/irisnet/irishub-sync/conf/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"net/http"
	"os"
)

var (
	logger Logger

	//zap method
	Binary = zap.Binary
	Bool   = zap.Bool
	//ByteString = zap.ByteString
	Complex128 = zap.Complex128
	Complex64  = zap.Complex64
	Float64    = zap.Float64
	Float32    = zap.Float32
	Int        = zap.Int
	Int64      = zap.Int64
	Int32      = zap.Int32
	Int16      = zap.Int16
	Int8       = zap.Int8
	String     = zap.String
	Uint       = zap.Uint
	Uint64     = zap.Uint64
	Uint32     = zap.Uint32
	Uint16     = zap.Uint16
	Uint8      = zap.Uint8
	Time       = zap.Time
	Any        = zap.Any
	Duration   = zap.Duration

	Info  = logger.Info
	Debug = logger.Debug
	Warn  = logger.Warn
	Error = logger.Error
	Panic = logger.Panic
	Fatal = logger.Fatal
	With  = logger.With
	Sync  = logger.Sync
)

type Logger struct {
	*zap.Logger
}

func init() {
	// 仅打印Info级别以上的日志
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})
	// 打印所有级别的日志
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})

	hook := lumberjack.Logger{
		Filename:   log.Conf.Filename,
		MaxSize:    log.Conf.MaxSize, // megabytes
		MaxBackups: 3,
		MaxAge:     log.Conf.MaxAge,   //days
		Compress:   log.Conf.Compress, // disabled by default
		LocalTime:  true,
	}

	fileWriter := zapcore.AddSync(&hook)

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	consoleDebugging := zapcore.Lock(os.Stdout)

	// Optimize the Kafka output for machine consumption and the console output
	// for human operators.
	consoleEncoder := zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig())

	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the four cores together.
	level := zap.NewAtomicLevelAt(zap.InfoLevel)
	var core zapcore.Core
	if log.Conf.EnableAtomicLevel {
		core = zapcore.NewTee(
			// 打印在控制台
			zapcore.NewCore(consoleEncoder, consoleDebugging, level),
			// 打印在文件中
			zapcore.NewCore(consoleEncoder, fileWriter, level),
		)
	} else {
		core = zapcore.NewTee(
			// 打印在控制台
			zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
			// 打印在文件中
			zapcore.NewCore(consoleEncoder, fileWriter, highPriority),
		)
	}

	caller := zap.AddCaller()
	callerSkip := zap.AddCallerSkip(1)
	// From a zapcore.Core, it's easy to construct a Logger.
	zapLogger := zap.New(core, caller, callerSkip)
	logger = Logger{
		zapLogger,
	}

	Info = logger.Info
	Debug = logger.Debug
	Warn = logger.Warn
	Error = logger.Error
	Panic = logger.Panic
	Fatal = logger.Fatal
	With = logger.With
	Sync = logger.Sync

	if log.Conf.EnableAtomicLevel {
		go func() {
			// curl -X PUT -H "Content-Type:application/json" -d '{"level":"info"}' localhost:9090
			http.ListenAndServe(":9090", &level)
		}()
	}
}
