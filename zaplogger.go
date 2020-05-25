package elog

import (
	"errors"
	"os"
	"sort"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	_ Logger = &zapLogger{}
)

type zapLogger struct {
	*zap.SugaredLogger
	level      zap.AtomicLevel
	encoding   string
	encoderCfg *EncoderConfig
}

func (l *zapLogger) SetLevel(level Level) {
	l.level.SetLevel(toZapLevel(level))
}

func (l *zapLogger) newZapOutput(name string, cfg *OutputConfig) (Output, *OutputConfig, zapcore.Level, error) {
	if cfg == nil {
		return nil, nil, zapcore.DebugLevel, errors.New("empty output configs")
	}
	if cfg.Encoding == "" {
		cfg.Encoding = l.encoding
	}
	if cfg.EncoderConfig == nil {
		cfg.EncoderConfig = l.encoderCfg
	} else {
		cfg.EncoderConfig.merge(l.encoderCfg)
	}
	out, err := newOutput(name, cfg)
	if err != nil {
		return nil, cfg, zapcore.DebugLevel, err
	}
	lvl := zapcore.DebugLevel
	if cfg.Level != "" {
		if _l, err := ParseLevel(cfg.Level); err == nil {
			lvl = toZapLevel(_l)
		}
	}
	return out, cfg, lvl, err
}

func (l *zapLogger) AddOutput(name string, cfg *OutputConfig) error {
	out, outCfg, outLevel, err := l.newZapOutput(name, cfg)
	if err != nil {
		return err
	}
	enc := l.newOutputEncoder(outCfg.Encoding, outCfg.EncoderConfig)
	newCore := &zapCore{
		minLevel: outLevel,
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

func (l *zapLogger) clone() *zapLogger {
	return &zapLogger{
		SugaredLogger: l.SugaredLogger,
		level:         l.level,
		encoding:      l.encoding,
		encoderCfg:    l.encoderCfg,
	}
}

func (l *zapLogger) Named(name string) Logger {
	nl := l.clone()
	nl.SugaredLogger = nl.SugaredLogger.Named(name)
	return nl
}

func (l *zapLogger) With(kvs ...interface{}) Logger {
	nl := l.clone()
	nl.SugaredLogger = l.SugaredLogger.With(kvs...)
	return nl
}

func (l *zapLogger) AddCallerSkip(skip int) Logger {
	clone := l.clone()
	unsugar := clone.SugaredLogger.Desugar()
	unsugar = unsugar.WithOptions(zap.AddCallerSkip(skip))
	clone.SugaredLogger = unsugar.Sugar()
	return clone
}

func (l *zapLogger) Printf(format string, args ...interface{}) {
	l.SugaredLogger.Infof(format, args...)
}
func (l *zapLogger) Println(args ...interface{}) {
	l.SugaredLogger.Info(args...)
}

func (l *zapLogger) newOutputEncoder(encoding string, cfg *EncoderConfig) zapcore.Encoder {
	if cfg == nil {
		return newZapEncoder(l.encoding, l.encoderCfg)
	}
	if encoding == "" {
		encoding = l.encoding
	}
	return newZapEncoder(encoding, cfg)
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

func buildZapOptions(cfg *Config) (zOpts []zap.Option) {
	// ErrOutputs
	zOpts = append(zOpts, zap.ErrorOutput(zapcore.Lock(os.Stderr)))
	if !cfg.CallerDisabled {
		zOpts = append(zOpts, zap.AddCaller())
	}
	if cfg.StacktraceEnabled {
		zOpts = append(zOpts, zap.AddStacktrace(zap.ErrorLevel))
	}
	zOpts = append(zOpts, zap.Fields(buildInitialFields(cfg.InitialFields)...))
	// Sampler
	// zOpts = append(zOpts, zap.WrapCore(func(core zapcore.Core) zapcore.Core {
	//   return zapcore.NewSampler(core, time.Second, 100, 100)
	// }))
	return
}

func New(cfg *Config) Logger {
	if cfg == nil {
		cfg = DefaultConfig()
	}
	cfg.validate()

	lg := &zapLogger{}
	zOpts := buildZapOptions(cfg)
	// level
	topLevel := zapcore.InfoLevel
	if lvl, err := ParseLevel(cfg.Level); err == nil {
		topLevel = toZapLevel(lvl)
	}
	lg.level = zap.NewAtomicLevelAt(topLevel)
	lg.encoderCfg = cfg.EncoderConfig
	lg.encoding = cfg.Encoding

	// Outputs
	cores := make([]zapcore.Core, 0)

	for name, _cfg := range cfg.Output {
		if !_cfg.Enabled {
			continue
		}
		out, outCfg, _level, err := lg.newZapOutput(name, _cfg)
		if err != nil {
			panic(err)
		}
		enc := lg.newOutputEncoder(outCfg.Encoding, outCfg.EncoderConfig)
		newCore := &zapCore{
			minLevel: _level,
			level:    lg.level,
			out:      out,
			enc:      enc,
		}
		cores = append(cores, newCore)
	}
	zl := zap.New(zapcore.NewTee(cores...), zOpts...)
	if cfg.Name != "" {
		zl = zl.Named(cfg.Name)
	}
	lg.SugaredLogger = zl.Sugar()
	return lg
}
