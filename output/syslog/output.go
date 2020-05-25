package syslog

import (
	"errors"
	"fmt"
	"log/syslog"
	"net/url"
	"os"
	"path/filepath"

	"github.com/derry6/elog"
	"github.com/derry6/elog/internal/param"
)

var (
	_ elog.Output = &Output{}
)

func init() {
	_ = elog.RegisterOutput(Name, New)
}

type Output struct {
	writer *syslog.Writer
}

func (o *Output) Close() error { return nil }
func (o *Output) Sync() error  { return nil }
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

func dial(cfg *elog.OutputConfig) (*syslog.Writer, error) {
	network := ""
	addr := ""
	params := param.Params(cfg.Params)
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

func New(cfg *elog.OutputConfig) (elog.Output, error) {
	if cfg == nil {
		return nil, errors.New("missing output configs")
	}
	w, err := dial(cfg)
	if err != nil {
		return nil, err
	}
	return &Output{writer: w}, nil
}
