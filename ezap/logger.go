package ezap

import (
	"os"
	"sort"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/derry6/elog"
	"github.com/derry6/elog/encoder"
	"github.com/derry6/elog/output"

	_ "github.com/derry6/elog/output/console"
	_ "github.com/derry6/elog/output/file"
	_ "github.com/derry6/elog/output/syslog"
)

var (
	_ elog.Logger = &logger{}
)

type logger struct {
	*zap.SugaredLogger
	level              zap.AtomicLevel
	defaultEncoding    string
	defaultEncoderOpts []encoder.Option
}

func (l *logger) SetLevel(levelStr string) {
	level := zapcore.DebugLevel
	if err := level.UnmarshalText([]byte(levelStr)); err == nil {
		l.level.SetLevel(level)
	}
}

func (l *logger) AddOutput(out output.Output) error {
	outOpts := out.Options()
	enc := l.newOutputEncoder(outOpts, outOpts.ColorEnabled)
	minLevel := zapcore.DebugLevel
	if lvl, err := elog.ParseLevel(outOpts.MinLevel); err == nil {
		minLevel = toZapLevel(lvl)
	}
	newCore := &zCore{
		minLevel: minLevel,
		out:      out,
		level:    l.level,
		enc:      enc,
	}
	unSugar := l.SugaredLogger.Desugar()
	wrapCore := zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewTee(core, newCore)
	})
	l.SugaredLogger = unSugar.WithOptions(wrapCore).Sugar()
	return nil
}

func (l *logger) clone() *logger {
	eOpts := make([]encoder.Option, len(l.defaultEncoderOpts))
	copy(eOpts, l.defaultEncoderOpts)
	return &logger{
		SugaredLogger:      l.SugaredLogger,
		level:              l.level,
		defaultEncoding:    l.defaultEncoding,
		defaultEncoderOpts: eOpts,
	}
}

func (l *logger) Named(name string) elog.Logger {
	nl := l.clone()
	nl.SugaredLogger = nl.SugaredLogger.Named(name)
	return nl
}

func (l *logger) With(kvs ...interface{}) elog.Logger {
	nl := l.clone()
	nl.SugaredLogger = l.SugaredLogger.With(kvs...)
	return nl
}

func (l *logger) AddCallerSkip(skip int) elog.Logger {
	clone := l.clone()
	unsugar := clone.SugaredLogger.Desugar()
	unsugar = unsugar.WithOptions(zap.AddCallerSkip(skip))
	clone.SugaredLogger = unsugar.Sugar()
	return clone
}

func (l *logger) Printf(format string, args ...interface{}) {
	l.SugaredLogger.Infof(format, args...)
}
func (l *logger) Println(args ...interface{}) {
	l.SugaredLogger.Info(args...)
}

func (l *logger) newOutputEncoder(outOpts *output.Options, colored bool) zapcore.Encoder {
	if outOpts == nil {
		return newEncoder(l.defaultEncoding, colored, l.defaultEncoderOpts...)
	}
	encoding := outOpts.Encoding
	encoderOpts := append(l.defaultEncoderOpts, outOpts.EncoderOpts...)
	if encoding == "" {
		encoding = l.defaultEncoding
	}
	return newEncoder(encoding, colored, encoderOpts...)
}

func buildInitialFields(initials map[string]interface{}) []zap.Field {
	if len(initials) == 0 {
		return []zap.Field{}
	}
	fs := make([]zap.Field, 0, len(initials))
	keys := make([]string, 0, len(initials))
	for k := range initials {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fs = append(fs, zap.Any(k, initials[k]))
	}
	return fs
}

func buildZapOptions(logOpts *elog.Options) (zOpts []zap.Option) {
	// ErrOutputs
	zOpts = append(zOpts, zap.ErrorOutput(zapcore.Lock(os.Stderr)))
	if !logOpts.CallerDisabled {
		zOpts = append(zOpts, zap.AddCaller())
	}
	if logOpts.StacktraceEnabled {
		zOpts = append(zOpts, zap.AddStacktrace(zap.ErrorLevel))
	}
	zOpts = append(zOpts, zap.Fields(buildInitialFields(logOpts.InitialFields)...))
	// Sampler
	// zOpts = append(zOpts, zap.WrapCore(func(core zapcore.Core) zapcore.Core {
	//   return zapcore.NewSampler(core, time.Second, 100, 100)
	// }))
	return
}

func New(opts ...elog.Option) elog.Logger {
	logOpts := elog.NewOptions(opts...)

	// fmt.Printf("===\n%s\n===\n", logOpts.Dump(true))

	lg := &logger{}
	zOpts := buildZapOptions(logOpts)
	// level
	zapLevel := zapcore.InfoLevel
	if lvl, err := elog.ParseLevel(logOpts.Level); err == nil {
		zapLevel = toZapLevel(lvl)
	}
	lg.level = zap.NewAtomicLevelAt(zapLevel)
	lg.defaultEncoding = logOpts.DefaultEncoding
	lg.defaultEncoderOpts = logOpts.DefaultEncoderOpts

	// Outputs
	cores := make([]zapcore.Core, 0)

	if len(logOpts.Outputs) == 0 && logOpts.ConsoleDisabled == false {
		logOpts.Outputs["console"] = []output.Option{}
	}
	for name, oOpts := range logOpts.Outputs {
		if logOpts.ConsoleDisabled && name == "console" {
			continue
		}
		out, err := output.New(name, oOpts...)
		if err != nil {
			panic(err)
		}
		outOpts := out.Options()
		minLevel := zapcore.DebugLevel
		if lvl, err := elog.ParseLevel(outOpts.MinLevel); err == nil {
			minLevel = toZapLevel(lvl)
		}
		colored := false
		if name == "console" && outOpts.ColorEnabled {
			colored = true
		}
		enc := lg.newOutputEncoder(outOpts, colored)
		newCore := &zCore{
			minLevel: minLevel,
			level:    lg.level,
			out:      out,
			enc:      enc,
		}
		cores = append(cores, newCore)
	}
	zl := zap.New(zapcore.NewTee(cores...), zOpts...)
	if logOpts.Name != "" {
		zl = zl.Named(logOpts.Name)
	}
	lg.SugaredLogger = zl.Sugar()
	return lg
}
