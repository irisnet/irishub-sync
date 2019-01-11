package logger

import "go.uber.org/zap"

type Logger struct {
	*zap.Logger
	module string
}

func GetLogger(module string) Logger {
	l := Logger{
		Logger: zapLogger,
		module: module,
	}
	return l
}

func (l Logger) appendModule(fields []Field) []Field {
	var f []Field
	if len(l.module) != 0 {
		f = append(f, String("module", l.module))
		f = append(f, fields...)
	}
	return f
}

func (l Logger) Info(msg string, fields ...Field) {
	defer Sync()
	Info(msg, l.appendModule(fields)...)
}

func (l Logger) Debug(msg string, fields ...Field) {
	defer Sync()
	Debug(msg, l.appendModule(fields)...)
}

func (l Logger) Warn(msg string, fields ...Field) {
	defer Sync()
	Warn(msg, l.appendModule(fields)...)
}

func (l Logger) Error(msg string, fields ...Field) {
	defer Sync()
	Error(msg, l.appendModule(fields)...)
}

func (l Logger) Panic(msg string, fields ...Field) {
	defer Sync()
	Panic(msg, l.appendModule(fields)...)
}

func (l Logger) Fatal(msg string, fields ...Field) {
	defer Sync()
	Fatal(msg, l.appendModule(fields)...)
}
