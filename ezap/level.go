package ezap

import (
    "go.uber.org/zap/zapcore"

    "github.com/derry6/elog"
)

func toZapLevel(lvl elog.Level) zapcore.Level {
    switch lvl {
    case elog.DEBUG:
        return zapcore.DebugLevel
    case elog.INFO:
        return zapcore.InfoLevel
    case elog.WARN:
        return zapcore.WarnLevel
    case elog.ERROR:
        return zapcore.ErrorLevel
    case elog.PANIC:
        return zapcore.PanicLevel
    case elog.FATAL:
        return zapcore.FatalLevel
    }
    return zapcore.Level(-1)
}

func fromZapLevel(lvl zapcore.Level) elog.Level {
    switch lvl {
    case zapcore.DebugLevel:
        return elog.DEBUG
    case zapcore.InfoLevel:
        return elog.INFO
    case zapcore.WarnLevel:
        return elog.WARN
    case zapcore.ErrorLevel:
        return elog.ERROR
    case zapcore.PanicLevel, zapcore.DPanicLevel:
        return elog.PANIC
    case zapcore.FatalLevel:
        return elog.FATAL
    }
    return elog.Level(-1)
}
