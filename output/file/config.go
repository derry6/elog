package file

import (
	"github.com/derry6/elog"
	"github.com/derry6/elog/internal/param"
)

const (
	Name           = "file"
	LogDir         = "logDir"
	BaseName       = "baseName"
	MultiFile      = "multiFile"
	RollingEnabled = "rollingEnabled"
	MaxSize        = "maxSize"
	MaxDays        = "maxDays"
	MaxBackups     = "maxBackups"
	Compressed     = "compressed"
	NameUseUTC     = "nameUseUTC"
)

func DefaultConfig() *elog.OutputConfig {
	return &elog.OutputConfig{
		Enabled: true,
		Params: map[string]interface{}{
			RollingEnabled: false,
		},
	}
}

type rollingConfig struct {
	maxSize    int
	maxAge     int
	maxBackups int
	compressed bool
	localTime  bool
}

func getRollingConfig(params param.Params) *rollingConfig {
	rollingEnabled := params.Bool(RollingEnabled, false)
	if rollingEnabled == false {
		return nil
	}
	var c rollingConfig
	c.maxSize = params.Int(MaxSize, 512)
	c.maxAge = params.Int(MaxDays, 7)
	c.maxBackups = params.Int(MaxBackups, 3)
	c.compressed = params.Bool(Compressed, false)
	c.localTime = !params.Bool(NameUseUTC, false)
	return &c
}
