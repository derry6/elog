package main

import (
	"github.com/derry6/elog"
	"github.com/derry6/elog/output/file"
)

func main() {
	elog.AddOutput(file.Name, file.DefaultConfig())
	msg := "This is info message This is info message This is info message " +
		"This is info message This is info message This is info message " +
		"This is info message This is info message"
	for i := 0; i < 10000000; i++ {
		elog.Infow(msg, "index", i)
	}
	_ = elog.Sync()
}
