package config

import (
	"github.com/derry6/elog"
	"github.com/derry6/elog/encoder"
)


type Config struct {
    Level             string                  `json:"level,omitempty" yaml:"level"`
    CallerDisabled    bool                    `json:"callerDisabled,omitempty" yaml:"callerDisabled"`
    StacktraceEnabled bool                    `json:"stacktraceEnabled,omitempty" yaml:"stacktraceEnabled"`
    Encoding          string                  `json:"encoding,omitempty" yaml:"encoding"`
    EncoderConfig     *encoder.Options        `json:"encoderConfig,omitempty" yaml:"encoderConfig"`
    InitialFields     map[string]interface{}  `json:"initialFields,omitempty" yaml:"initialFields"`
    Outputs           map[string]*OutputConfig `json:"outputs" yaml:"outputs"`
}

func DefaultConfig() *Config {
    return &Config{
        Level:             "INFO",
        CallerDisabled:    false,
        StacktraceEnabled: false,
        Encoding:          "text",
        EncoderConfig:     nil,
        InitialFields:     nil,
        Outputs:           nil,
    }
}

func (c *Config) BuildOptions() (opts []elog.Option) {
    addStringOpt := func(v string, withFn func(string) elog.Option) {
        if v != "" {
            opts = append(opts, withFn(v))
        }
    }
    addStringOpt(c.Level, elog.WithLevel)
    if c.CallerDisabled {
        opts = append(opts, elog.WithCallerDisabled())
    }
    if c.StacktraceEnabled {
        opts = append(opts, elog.WithStacktraceEnabled())
    }
    addStringOpt(c.Encoding, elog.WithDefaultEncoding)
    if eOpts := BuildEncoderOptions(c.EncoderConfig); len(eOpts) > 0 {
        opts = append(opts, elog.WithDefaultEncoderOptions(eOpts...))
    }
    if len(c.InitialFields) > 0 {
        opts = append(opts, elog.WithInitialFields(c.InitialFields))
    }
    for outputName, outputConfig := range c.Outputs {
        if outputConfig.Enabled && outputName != "" {
            outOpts := outputConfig.BuildOptions()
            opts = append(opts, elog.WithOutput(outputName, outOpts...))
        }
    }
    return opts
}
