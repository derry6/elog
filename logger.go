package elog

// Logger logger.yaml interface
type Logger interface {
	SetLevel(levelStr Level)
	Sync() error

	// new loggers
	Named(name string) Logger
	With(kvs ...interface{}) Logger
	AddCallerSkip(skip int) Logger

	AddOutput(name string, cfg *OutputConfig) error

	// elog writers
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})

	// elog key values
	Debugw(msg string, kvs ...interface{})
	Infow(msg string, kvs ...interface{})
	Warnw(msg string, kvs ...interface{})
	Errorw(msg string, kvs ...interface{})
	Fatalw(msg string, kvs ...interface{})
	Panicw(msg string, kvs ...interface{})

	// printf-like
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Printf(format string, args ...interface{})
	Println(args ...interface{})
}
