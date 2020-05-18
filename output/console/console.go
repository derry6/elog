package console

import (
	"io"
	"io/ioutil"
	"os"
	"sync"

	"github.com/derry6/elog/internal/param"
	"github.com/derry6/elog/internal/terminal"
	"github.com/derry6/elog/output"
)

const (
	Name    = "console"
	Writer  = "writer"
	Colored = "colored"
)

func init() {
	_ = output.Register(Name, New)
}

type consoleOut struct {
	opts *output.Options
	w    io.Writer
	mu   sync.Mutex
}

func (c *consoleOut) Close() error             { return nil }
func (c *consoleOut) Options() *output.Options { return c.opts }
func (c *consoleOut) Write(_ output.Level, line []byte) error {
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

func New(opts ...output.Option) (output.Output, error) {
	outOpts := output.NewOptions(opts...)
	w := io.Writer(os.Stderr)
	p := param.Params(outOpts.Params)
	switch p.String(Writer, "stderr") {
	case "stdout":
		w = os.Stdout
	case "discard":
		w = ioutil.Discard
	}
	if p.Exists(Colored) {
		outOpts.ColorEnabled = p.Bool(Colored, true)
	} else {
		if terminal.Check(w) {
			outOpts.ColorEnabled = true
		}
	}
	outOpts.Params = map[string]interface{}{}
	out := &consoleOut{w: w, opts: outOpts}
	return out, nil
}
