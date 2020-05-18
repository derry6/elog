package ezap

import (
    "go.uber.org/zap/zapcore"

    "github.com/derry6/elog/output"
)

type zCore struct {
	level    zapcore.LevelEnabler
	minLevel zapcore.LevelEnabler
	out      output.Output
	enc      zapcore.Encoder
}

func (c *zCore) Enabled(l zapcore.Level) bool {
	return c.level.Enabled(l) && c.minLevel.Enabled(l)
}

func (c *zCore) With(fields []zapcore.Field) zapcore.Core {
	clone := c.clone()
	for i := range fields {
		fields[i].AddTo(c.enc)
	}
	return clone
}

func (c *zCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c *zCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	buf, err := c.enc.EncodeEntry(ent, fields)
	if err != nil {
		return err
	}
	err = c.out.Write(fromZapLevel(ent.Level), buf.Bytes())
	buf.Free()
	if err != nil {
		return err
	}
	if ent.Level > zapcore.ErrorLevel {
		_ = c.Sync()
	}
	return nil
}

func (c *zCore) Sync() error {
	return c.out.Sync()
}

func (c *zCore) clone() *zCore {
	return &zCore{
		level:    c.level,
		minLevel: c.minLevel,
		enc:      c.enc.Clone(),
		out:      c.out,
	}
}
