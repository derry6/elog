package encoder

type Options struct {
	MessageKey      string `json:"messageKey,omitempty" yaml:"messageKey"`
	LevelKey        string `json:"levelKey,omitempty" yaml:"levelKey"`
	TimeKey         string `json:"timeKey,omitempty" yaml:"timeKey"`
	NameKey         string `json:"nameKey,omitempty" yaml:"nameKey"`
	CallerKey       string `json:"callerKey,omitempty" yaml:"callerKey"`
	StacktraceKey   string `json:"stacktraceKey,omitempty" yaml:"stacktraceKey"`
	LineEnding      string `json:"lineEnding,omitempty" yaml:"lineEnding"`
	LevelEncoder    string `json:"levelEncoder,omitempty" yaml:"levelEncoder"`
	TimeEncoder     string `json:"timeEncoder,omitempty" yaml:"timeEncoder"`
	DurationEncoder string `json:"durationEncoder,omitempty" yaml:"durationEncoder"`
	CallerEncoder   string `json:"callerEncoder,omitempty" yaml:"callerEncoder"`
	NameEncoder     string `json:"nameEncoder,omitempty" yaml:"nameEncoder"`
}

type Option func(opts *Options)

func defaultEncoderOptions() *Options {
	return &Options{
		TimeKey:         "ts",
		LevelKey:        "level",
		NameKey:         "logger",
		CallerKey:       "caller",
		MessageKey:      "msg",
		StacktraceKey:   "stacktrace",
		LineEnding:      "\n",
		LevelEncoder:    "lowercase",
		TimeEncoder:     "rfc3339",
		DurationEncoder: "seconds",
		CallerEncoder:   "short",
	}
}

func WithMessageKey(name string) Option {
	return func(opts *Options) {
		opts.MessageKey = name
	}
}
func WithLevelKey(name string) Option {
	return func(opts *Options) {
		opts.LevelKey = name
	}
}
func WithTimeKey(name string) Option {
	return func(opts *Options) {
		opts.TimeKey = name
	}
}
func WithNameKey(name string) Option {
	return func(opts *Options) {
		opts.NameKey = name
	}
}
func WithCallerKey(name string) Option {
	return func(opts *Options) {
		opts.CallerKey = name
	}
}
func WithStacktraceKey(name string) Option {
	return func(opts *Options) {
		opts.StacktraceKey = name
	}
}

func WithLineEnding(lineEnding string) Option {
	return func(opts *Options) {
		opts.LineEnding = lineEnding
	}
}

func WithLevelEncoder(name string) Option {
	return func(opts *Options) {
		opts.LevelEncoder = name
	}
}
func WithDurationEncoder(name string) Option {
	return func(opts *Options) {
		opts.DurationEncoder = name
	}
}

func WithTimeEncoder(name string) Option {
	return func(opts *Options) {
		opts.TimeEncoder = name
	}
}
func WithCallerEncoder(name string) Option {
	return func(opts *Options) {
		opts.CallerEncoder = name
	}
}
func WithNameEncoder(name string) Option {
	return func(opts *Options) {
		opts.NameEncoder = name
	}
}

func NewOptions(opts ...Option) *Options {
	eOpts := defaultEncoderOptions()
	for _, o := range opts {
		o(eOpts)
	}
	return eOpts
}
