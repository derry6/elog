package main

import (
	"github.com/derry6/elog"
	"github.com/derry6/elog/ezap"
	"github.com/derry6/elog/output"
)

func main() {
	elog.Use(ezap.New(
		elog.WithFile("rolling",
			output.WithParam("rollingEnabled", true))))
	msg := "This is info message This is info message This is info message " +
		"This is info message This is info message This is info message " +
		"This is info message This is info message"
	for i := 0; i < 10000000; i++ {
		elog.Infow(msg, "index", i)
	}
	_ = elog.Sync()
}
