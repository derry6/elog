package elog

var (
	l0 Logger
	l1 Logger
)

func init() {
	l0 = New(DefaultConfig())
	_ = l0.AddOutput(Console, DefaultConsoleConfig())
	l1 = l0.AddCallerSkip(1)
}

// Get gets global logger
func Get() Logger { return l0 }

// Use sets global logger
func Use(logger Logger) {
	if logger == nil {
		return
	}
	l0 = logger
	l1 = l0.AddCallerSkip(1)
}

// Sync flush logs
func Sync() error {
	return l1.Sync()
}

func SetLevel(level Level)           { l1.SetLevel(level) }
func Named(name string) Logger       { return l1.Named(name) }
func With(kvs ...interface{}) Logger { return l1.With(kvs...) }
func AddCallerSkip(skip int) Logger  { return l1.AddCallerSkip(1) }
func AddOutput(name string, cfg *OutputConfig) error {
	err := l0.AddOutput(name, cfg)
	if err != nil {
		return err
	}
	l1 = l0.AddCallerSkip(1)
	return err
}

// elog writers
func Debug(args ...interface{}) { l1.Debug(args...) }
func Info(args ...interface{})  { l1.Info(args...) }
func Warn(args ...interface{})  { l1.Warn(args...) }
func Error(args ...interface{}) { l1.Error(args...) }
func Fatal(args ...interface{}) { l1.Fatal(args...) }
func Panic(args ...interface{}) { l1.Panic(args...) }

// elog key values
func Debugw(msg string, kvs ...interface{}) { l1.Debugw(msg, kvs...) }
func Infow(msg string, kvs ...interface{})  { l1.Infow(msg, kvs...) }
func Warnw(msg string, kvs ...interface{})  { l1.Warnw(msg, kvs...) }
func Errorw(msg string, kvs ...interface{}) { l1.Errorw(msg, kvs...) }
func Fatalw(msg string, kvs ...interface{}) { l1.Fatalw(msg, kvs...) }
func Panicw(msg string, kvs ...interface{}) { l1.Panicw(msg, kvs...) }

// printf-like
func Debugf(format string, args ...interface{}) { l1.Debugf(format, args...) }
func Infof(format string, args ...interface{})  { l1.Infof(format, args...) }
func Warnf(format string, args ...interface{})  { l1.Warnf(format, args...) }
func Errorf(format string, args ...interface{}) { l1.Errorf(format, args...) }
func Fatalf(format string, args ...interface{}) { l1.Fatalf(format, args...) }
func Panicf(format string, args ...interface{}) { l1.Panicf(format, args...) }

func Printf(format string, args ...interface{}) { l1.Printf(format, args...) }
func Println(args ...interface{})               { l1.Println(args...) }
