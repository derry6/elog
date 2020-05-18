package ezap

import (
    "strings"
    "time"

    "go.uber.org/zap/zapcore"

    "github.com/derry6/elog/encoder"
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

func newTimeEncoder(name string) zapcore.TimeEncoder {
    encoder := zapcore.TimeEncoder(zapcore.ISO8601TimeEncoder)
    if name != "" {
        name = strings.ToLower(name)
        format, ok := timeFormats[name]
        if ok && format != "" {
            return func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
                encoder.AppendString(time.Format(format))
            }
        }
        _ = encoder.UnmarshalText([]byte(name))
    }
    return encoder
}
func newLevelEncoder(levelFormat string, colored bool) zapcore.LevelEncoder {
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

func newEncoder(encoding string, colored bool, opts ...encoder.Option) zapcore.Encoder {
    eOpts := encoder.NewOptions(opts...)
    zEc := zapcore.EncoderConfig{
        TimeKey:        eOpts.TimeKey,
        LevelKey:       eOpts.LevelKey,
        NameKey:        eOpts.NameKey,
        CallerKey:      eOpts.CallerKey,
        MessageKey:     eOpts.MessageKey,
        StacktraceKey:  eOpts.StacktraceKey,
        LineEnding:     eOpts.LineEnding,
        EncodeLevel:    newLevelEncoder(eOpts.LevelEncoder, colored),
        EncodeTime:     newTimeEncoder(eOpts.TimeEncoder),
        EncodeDuration: zapcore.SecondsDurationEncoder,
        EncodeCaller:   zapcore.ShortCallerEncoder,
    }
    _ = zEc.EncodeDuration.UnmarshalText([]byte(eOpts.DurationEncoder))
    _ = zEc.EncodeName.UnmarshalText([]byte(eOpts.NameEncoder))
    switch encoding {
    case encoder.Json:
        return zapcore.NewJSONEncoder(zEc)
    default:
        return zapcore.NewConsoleEncoder(zEc)
    }
}
