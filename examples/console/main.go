package main

import (
	"github.com/derry6/elog"
)

func main() {
	elog.SetLevel(elog.DEBUG)
	elog.Debugw("This is debug message")
	elog.Infow("This is info message")
	elog.Errorw("This is error message")
}
