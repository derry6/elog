package config

import (
	"github.com/derry6/elog/encoder"
	"github.com/derry6/elog/internal/param"
	"github.com/derry6/elog/output"
)

type OutputConfig struct {
	Enabled       bool             `json:"enabled,omitempty" yaml:"enabled"`
	Level         string           `json:"level,omitempty" yaml:"level"`
	Encoding      string           `json:"encoding,omitempty" yaml:"encoding"`
	EncoderConfig *encoder.Options `json:"encoderConfig,omitempty" yaml:"encoderConfig"`
	Params        param.Params     `json:"parameters,omitempty" yaml:"parameters"`
}

func (c *OutputConfig) BuildOptions() (outOpts []output.Option) {
	addStrOpt := func(v string, withFn func(string) output.Option) {
		if v != "" {
			outOpts = append(outOpts, withFn(v))
		}
	}
	addStrOpt(c.Level, output.WithMinLevel)
	addStrOpt(c.Encoding, output.WithEncoding)

	if eOpts := BuildEncoderOptions(c.EncoderConfig); len(eOpts) > 0 {
		outOpts = append(outOpts, output.WithEncoderOptions(eOpts...))
	}
	if len(c.Params) > 0 {
		outOpts = append(outOpts, output.WithParams(c.Params))
	}
	return outOpts
}
