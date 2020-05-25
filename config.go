package elog

const (
	JSON = "json"
	Text = "text"
)

type EncoderConfig struct {
	ColorEnabled    bool   `json:"colorEnabled" yaml:"colorEnabled"`
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

func (c1 *EncoderConfig) merge(c2 *EncoderConfig) {
	if c2 == nil {
		return
	}
	_Ms := func(v1 *string, v2 string) {
		if *v1 == "" && v2 != "" {
			*v1 = v2
		}
	}
	_Ms(&c1.MessageKey, c2.MessageKey)
	_Ms(&c1.LevelKey, c2.LevelKey)
	_Ms(&c1.TimeKey, c2.TimeKey)
	_Ms(&c1.NameKey, c2.NameKey)
	_Ms(&c1.CallerKey, c2.CallerKey)
	_Ms(&c1.StacktraceKey, c2.StacktraceKey)
	_Ms(&c1.LineEnding, c2.LineEnding)
	_Ms(&c1.LevelEncoder, c2.LevelEncoder)
	_Ms(&c1.TimeEncoder, c2.TimeEncoder)
	_Ms(&c1.DurationEncoder, c2.DurationEncoder)
	_Ms(&c1.CallerEncoder, c2.CallerEncoder)
	_Ms(&c1.NameEncoder, c2.NameEncoder)
}

func DefaultEncoderConfig() *EncoderConfig {
	return &EncoderConfig{
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

type OutputConfig struct {
	Enabled       bool                   `json:"enabled,omitempty" yaml:"enabled"`
	Level         string                 `json:"level,omitempty" yaml:"level"`
	Encoding      string                 `json:"encoding,omitempty" yaml:"encoding"`
	EncoderConfig *EncoderConfig         `json:"encoderConfig,omitempty" yaml:"encoderConfig"`
	Params        map[string]interface{} `json:"parameters,omitempty" yaml:"parameters"`
}

type Config struct {
	Level             string                   `json:"level,omitempty" yaml:"level"`
	Name              string                   `json:"name" yaml:"name"`
	CallerDisabled    bool                     `json:"callerDisabled,omitempty" yaml:"callerDisabled"`
	StacktraceEnabled bool                     `json:"stacktraceEnabled,omitempty" yaml:"stacktraceEnabled"`
	Encoding          string                   `json:"encoding,omitempty" yaml:"encoding"`
	EncoderConfig     *EncoderConfig           `json:"encoderConfig,omitempty" yaml:"encoderConfig"`
	InitialFields     map[string]interface{}   `json:"initialFields,omitempty" yaml:"initialFields"`
	Output            map[string]*OutputConfig `json:"output" yaml:"output"`
}

func (c *Config) validate() {
	if c.Level == "" {
		c.Level = "INFO"
	}
	if c.Encoding == "" {
		c.Encoding = JSON
	}
	if c.EncoderConfig == nil {
		c.EncoderConfig = DefaultEncoderConfig()
	}
}

func DefaultConfig() *Config {
	return &Config{
		Level:             "INFO",
		CallerDisabled:    false,
		StacktraceEnabled: false,
		Encoding:          Text,
		EncoderConfig:     nil,
		InitialFields:     nil,
		Output:            nil,
	}
}
