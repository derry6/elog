package elog

import (
	"go.uber.org/zap/zapcore"
)

type zapCore struct {
	level    zapcore.LevelEnabler
	minLevel zapcore.LevelEnabler
	enc      zapcore.Encoder
	out      Output
}

func (c *zapCore) Enabled(l zapcore.Level) bool {
	return c.level.Enabled(l) && c.minLevel.Enabled(l)
}

func (c *zapCore) With(fields []zapcore.Field) zapcore.Core {
	clone := c.clone()
	for i := range fields {
		fields[i].AddTo(c.enc)
	}
	return clone
}

func (c *zapCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c *zapCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
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

func (c *zapCore) Sync() error {
	return c.out.Sync()
}

func (c *zapCore) clone() *zapCore {
	return &zapCore{
		level:    c.level,
		minLevel: c.minLevel,
		enc:      c.enc.Clone(),
		out:      c.out,
	}
}
