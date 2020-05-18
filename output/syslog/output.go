package syslog

import (
	"fmt"
	"log/syslog"
	"net/url"
	"os"
	"path/filepath"

	"github.com/derry6/elog"
	"github.com/derry6/elog/internal/param"
	"github.com/derry6/elog/output"
)

const (
	Name     = "syslog"
	Host     = "host"
	Facility = "facility"
	Tag      = "tag"
)

var (
	_ output.Output = &Output{}
)

func init() {
	_ = output.Register(Name, New)
}

type Output struct {
	opts   *output.Options
	writer *syslog.Writer
}

func (o *Output) Close() error             { return nil }
func (o *Output) Sync() error              { return nil }
func (o *Output) Options() *output.Options { return o.opts }
func (o *Output) Write(l elog.Level, line []byte) error {
	switch l {
	case elog.DEBUG:
		return o.writer.Debug(bytes2Str(line))
	case elog.INFO:
		return o.writer.Info(bytes2Str(line))
	case elog.WARN:
		return o.writer.Warning(bytes2Str(line))
	case elog.ERROR:
		return o.writer.Err(bytes2Str(line))
	case elog.PANIC, elog.FATAL:
		return o.writer.Emerg(bytes2Str(line))
	default:
		return o.writer.Info(bytes2Str(line))
	}
}

func dial(outOpts *output.Options) (*syslog.Writer, error) {
	network := ""
	addr := ""
	params := param.Params(outOpts.Params)
	fac := params.String(Facility, "")
	f, err := parseFacility(fac)
	if err != nil {
		return nil, err
	}
	host := params.String(Host, "")
	if host == "" {
		host = "localhost"
	} else if host != "localhost" {
		u, err := url.Parse(host)
		if err != nil {
			return nil, err
		}
		network = u.Scheme
		p := u.Port()
		if p == "" {
			p = "514"
		}
		addr = fmt.Sprintf("%s:%s", u.Host, p)
	}
	tag := params.String(Tag, "")
	if tag == "" {
		tag = filepath.Base(os.Args[0])
	}
	return syslog.Dial(network, addr, f, tag)
}

func New(opts ...output.Option) (output.Output, error) {
	outOpts := output.NewOptions(opts...)
	w, err := dial(outOpts)
	if err != nil {
		return nil, err
	}
	return &Output{writer: w, opts: outOpts}, nil
}
