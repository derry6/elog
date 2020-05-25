package elog

import (
	"strings"
	"time"

	"go.uber.org/zap/zapcore"
)

var (
	timeFormats = map[string]string{
		"ansic":       time.ANSIC,
		"unixdate":    time.UnixDate,
		"rubydate":    time.RubyDate,
		"rfc822":      time.RFC822,
		"rfc822z":     time.RFC822Z,
		"rfc850":      time.RFC850,
		"rfc1123":     time.RFC1123,
		"rfc1123z":    time.RFC1123Z,
		"rfc3339":     time.RFC3339,
		"rfc3339nano": time.RFC3339Nano,
		"kitchen":     time.Kitchen,
		"stamp":       time.Stamp,
		"stampmicro":  time.StampMicro,
		"stampmilli":  time.StampMilli,
		"stampnano":   time.StampNano,
		"iso8601":     "2006-01-02T15:04:05.000Z0700",
	}
)

func newZapTimeEncoder(name string) zapcore.TimeEncoder {
	encoder := zapcore.TimeEncoder(zapcore.ISO8601TimeEncoder)
	if name != "" {
		name = strings.ToLower(name)
		format, ok := timeFormats[name]
		if ok && format != "" {
			return func(time time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(time.Format(format))
			}
		}
		_ = encoder.UnmarshalText([]byte(name))
	}
	return encoder
}

func newZapLevelEncoder(levelFormat string, colored bool) zapcore.LevelEncoder {
	if len(levelFormat) == 0 {
		if colored {
			return zapcore.CapitalColorLevelEncoder
		}
		return zapcore.CapitalLevelEncoder
	}
	switch levelFormat {
	case "upper", "uppercase", "capital":
		if colored {
			return zapcore.CapitalColorLevelEncoder
		}
		return zapcore.CapitalLevelEncoder
	default:
		if colored {
			return zapcore.LowercaseColorLevelEncoder
		}
		return zapcore.LowercaseLevelEncoder
	}
}

func newZapEncoder(encoding string, cfg *EncoderConfig) zapcore.Encoder {
	if cfg == nil {
		cfg = DefaultEncoderConfig()
	}
	zEc := zapcore.EncoderConfig{
		TimeKey:        cfg.TimeKey,
		LevelKey:       cfg.LevelKey,
		NameKey:        cfg.NameKey,
		CallerKey:      cfg.CallerKey,
		MessageKey:     cfg.MessageKey,
		StacktraceKey:  cfg.StacktraceKey,
		LineEnding:     cfg.LineEnding,
		EncodeLevel:    newZapLevelEncoder(cfg.LevelEncoder, cfg.ColorEnabled),
		EncodeTime:     newZapTimeEncoder(cfg.TimeEncoder),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	_ = zEc.EncodeDuration.UnmarshalText([]byte(cfg.DurationEncoder))
	_ = zEc.EncodeName.UnmarshalText([]byte(cfg.NameEncoder))
	switch encoding {
	case JSON:
		return zapcore.NewJSONEncoder(zEc)
	default:
		return zapcore.NewConsoleEncoder(zEc)
	}
}
