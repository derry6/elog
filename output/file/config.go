package file

import (
	"github.com/derry6/elog/internal/param"
)

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
