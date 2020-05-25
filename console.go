package elog

import (
	"io"
	"io/ioutil"
	"os"
	"sync"
)

const (
	Console         = "console"
	_consoleWriter  = "writer"
	_consoleColored = "colored"
)

func DefaultConsoleConfig() *OutputConfig {
	return &OutputConfig{
		Enabled: true,
		Params: map[string]interface{}{
			_consoleWriter: "stderr",
		},
	}
}

func init() {
	_ = RegisterOutput(Console, newConsole)
}

type consoleOut struct {
	w  io.Writer
	mu sync.Mutex
}

func (c *consoleOut) Close() error { return nil }
func (c *consoleOut) Write(_ Level, line []byte) error {
	c.mu.Lock()
	_, err := c.w.Write(line)
	c.mu.Unlock()
	return err
}
func (c *consoleOut) Sync() error {
	c.mu.Lock()
	if f, ok := c.w.(*os.File); ok && f != nil {
		c.mu.Unlock()
		return f.Sync()
	}
	c.mu.Lock()
	return nil
}

func newConsole(cfg *OutputConfig) (Output, error) {
	// w := io.Writer(os.Stderr)
	out := &consoleOut{w: os.Stderr}
	if cfg == nil {
		return out, nil
	}
	writer, ok := cfg.Params[_consoleWriter]
	if ok {
		switch w := writer.(type) {
		case string:
			if w == "stdout" {
				out.w = os.Stdout
			}
			if w == "discard" {
				out.w = ioutil.Discard
			}
		}
	}
	return out, nil
}
