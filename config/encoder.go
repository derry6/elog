package config

import "github.com/derry6/elog/encoder"

func BuildEncoderOptions(c *encoder.Options) (opts []encoder.Option) {
    if c == nil {
        return
    }
    addOpt := func(name string, withFn func(string) encoder.Option) {
        if name != "" {
            opts = append(opts, withFn(name))
        }
    }
    addOpt(c.LineEnding, encoder.WithLineEnding)
    addOpt(c.MessageKey, encoder.WithMessageKey)
    addOpt(c.LevelKey, encoder.WithLevelKey)
    addOpt(c.TimeKey, encoder.WithTimeKey)
    addOpt(c.NameKey, encoder.WithNameKey)
    addOpt(c.CallerKey, encoder.WithCallerKey)
    addOpt(c.StacktraceKey, encoder.WithStacktraceKey)
    addOpt(c.LineEnding, encoder.WithLineEnding)
    addOpt(c.DurationEncoder, encoder.WithDurationEncoder)
    addOpt(c.LevelEncoder, encoder.WithLevelEncoder)
    addOpt(c.TimeEncoder, encoder.WithTimeEncoder)
    addOpt(c.CallerEncoder, encoder.WithCallerEncoder)
    addOpt(c.NameEncoder, encoder.WithNameEncoder)
    return opts
}
