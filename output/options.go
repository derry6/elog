package output

import (
	"encoding/json"
	"github.com/derry6/elog/encoder"
)

type Options struct {
	MinLevel     string
	Encoding     string
	EncoderOpts  []encoder.Option
	ColorEnabled bool
	Params       map[string]interface{}
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

func defaultOutputOptions() *Options {
	return &Options{
		MinLevel:    "DEBUG",
		Encoding:    "",
		EncoderOpts: nil,
		Params:      map[string]interface{}{},
	}
}

func WithMinLevel(minLevel string) Option {
	return func(opts *Options) {
		opts.MinLevel = minLevel
	}
}
func WithEncoding(encoding string) Option {
	return func(opts *Options) {
		opts.Encoding = encoding
	}
}

func WithEncoderOptions(eOpts ...encoder.Option) Option {
	return func(opts *Options) {
		opts.EncoderOpts = eOpts
	}
}

func WithParam(k string, v interface{}) Option {
	return func(opts *Options) {
		if opts.Params == nil {
			opts.Params = map[string]interface{}{}
		}
		opts.Params[k] = v
	}
}

func WithParams(values map[string]interface{}) Option {
	return func(opts *Options) {
		if opts.Params == nil {
			opts.Params = map[string]interface{}{}
		}
		for k, v := range values {
			opts.Params[k] = v
		}
	}
}

type Option func(opts *Options)

func NewOptions(opts ...Option) *Options {
	outOpts := defaultOutputOptions()
	for _, o := range opts {
		o(outOpts)
	}
	return outOpts
}
