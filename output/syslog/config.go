package syslog

import "github.com/derry6/elog"

const (
	Name     = "syslog"
	Host     = "host"
	Facility = "facility"
	Tag      = "tag"
)

func DefaultConfig() *elog.OutputConfig {
	return &elog.OutputConfig{
		Enabled: true,
		Params: map[string]interface{}{
			Host: "localhost",
		},
	}
}
