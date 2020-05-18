package elog

import (
    "github.com/derry6/elog/output"
)

type Level = output.Level

const (
    DEBUG = output.DEBUG
    INFO  = output.INFO
    WARN  = output.WARN
    ERROR = output.ERROR
    PANIC = output.PANIC
    FATAL = output.FATAL
)

var (
    ParseLevel = output.ParseLevel
)
