package main

import (
    "github.com/derry6/elog"
    "github.com/derry6/elog/ezap"
)

func main() {
    elog.Use(ezap.New(elog.WithCallerDisabled()))
    elog.Debugw("This is debug message")
    elog.Infow("This is info message")
    elog.Errorw("This is error message")
}
