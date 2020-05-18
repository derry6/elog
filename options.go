package elog

import (
	"encoding/json"
	"github.com/derry6/elog/output/console"
	"github.com/derry6/elog/output/file"

	"github.com/derry6/elog/encoder"
	"github.com/derry6/elog/output"
)

// Logger

type Options struct {
	Name              string
	Level             string
	ConsoleDisabled   bool
	CallerDisabled    bool
	StacktraceEnabled bool
	// Defaults
	DefaultEncoding    string
	DefaultEncoderOpts []encoder.Option
	// Initials
	InitialFields map[string]interface{}
	// Outputs
	Outputs map[string][]output.Option
}

func (o *Options) Dump(pretty ...bool) string {
	var dat []byte
	if len(pretty) > 0 && pretty[0] {
		dat, _ = json.MarshalIndent(o, "", "  ")
	} else {
		dat, _ = json.Marshal(o)
	}
	return string(dat)
}

func defaultOptions() *Options {
	return &Options{
		Name:               "",
		Level:              "INFO",
		ConsoleDisabled:    false,
		CallerDisabled:     false,
		StacktraceEnabled:  false,
		DefaultEncoding:    "text",
		DefaultEncoderOpts: nil,
		InitialFields:      map[string]interface{}{},
		Outputs:            map[string][]output.Option{},
	}
}

type Option func(opts *Options)

func WithName(name string) Option {
	return func(opts *Options) {
		opts.Name = name
	}
}

func WithLevel(level string) Option {
	return func(opts *Options) {
		opts.Level = level
	}
}

func WithConsoleDisabled() Option {
	return func(opts *Options) {
		opts.ConsoleDisabled = true
	}
}

func WithCallerDisabled() Option {
	return func(opts *Options) {
		opts.CallerDisabled = true
	}
}

func WithStacktraceEnabled() Option {
	return func(opts *Options) {
		opts.StacktraceEnabled = true
	}
}

func WithDefaultEncoding(encoding string) Option {
	return func(opts *Options) {
		opts.DefaultEncoding = encoding
	}
}

func WithDefaultEncoderOptions(eOpts ...encoder.Option) Option {
	return func(opts *Options) {
		opts.DefaultEncoderOpts = eOpts
	}
}

func WithInitialFields(fields map[string]interface{}) Option {
	return func(opts *Options) {
		opts.InitialFields = fields
	}
}

func WithConsole(name string, outOpts ...output.Option) Option {
	return func(opts *Options) {
		opts.Outputs["console"] = append(outOpts,
			output.WithParam(console.Writer, name))
	}
}

func WithFile(baseName string, outOpts ...output.Option) Option {
	return func(opts *Options) {
		opts.Outputs["file"] = append(outOpts,
			output.WithParam(file.BaseName, baseName))
	}
}

func WithOutput(name string, outOpts ...output.Option) Option {
	return func(opts *Options) {
		opts.Outputs[name] = outOpts
	}
}

func NewOptions(opts ...Option) *Options {
	lOpts := defaultOptions()
	for _, o := range opts {
		o(lOpts)
	}
	return lOpts
}
