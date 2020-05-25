package elog

import (
	"go.uber.org/zap/zapcore"
)

func toZapLevel(lvl Level) zapcore.Level {
	switch lvl {
	case DEBUG:
		return zapcore.DebugLevel
	case INFO:
		return zapcore.InfoLevel
	case WARN:
		return zapcore.WarnLevel
	case ERROR:
		return zapcore.ErrorLevel
	case PANIC:
		return zapcore.PanicLevel
	case FATAL:
		return zapcore.FatalLevel
	}
	return zapcore.Level(-1)
}

func fromZapLevel(lvl zapcore.Level) Level {
	switch lvl {
	case zapcore.DebugLevel:
		return DEBUG
	case zapcore.InfoLevel:
		return INFO
	case zapcore.WarnLevel:
		return WARN
	case zapcore.ErrorLevel:
		return ERROR
	case zapcore.PanicLevel, zapcore.DPanicLevel:
		return PANIC
	case zapcore.FatalLevel:
		return FATAL
	}
	return Level(-1)
}
