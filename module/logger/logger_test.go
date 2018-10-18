package logger

import (
	"testing"
)

func TestZap(t *testing.T) {
	Info("this is a info logger", Int("test", 1))
	Debug("this is a debug logger")
	Warn("this is a warn logger")
	Error("this is a a error logger", Int("test", 2))
	With(Int("test", 2)).Warn("this is a warn logger")
	Sync()
}
