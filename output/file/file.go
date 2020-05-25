package file

import (
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/derry6/elog"
	"github.com/derry6/elog/internal/param"

	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	_ = elog.RegisterOutput(Name, New)
}

type fileWriter struct {
	l elog.Level
	w *lumberjack.Logger
}
type fileOutput struct {
	lws []fileWriter
}

func (o *fileOutput) Close() error {
	for _, l := range o.lws {
		_ = l.w.Close()
	}
	return nil
}
func (o *fileOutput) Sync() error {
	return nil
}
func (o *fileOutput) Write(l elog.Level, line []byte) (err error) {
	for _, x := range o.lws {
		if l >= x.l {
			_, err = x.w.Write(line)
			return
		}
	}
	return err
}

func newFileLogger(filename string, c *rollingConfig) *lumberjack.Logger {
	if c == nil {
		return &lumberjack.Logger{
			Filename:   filename,
			MaxSize:    math.MaxInt32 / (1024 * 1024),
			MaxAge:     0,
			MaxBackups: 0,
			Compress:   false,
			LocalTime:  true,
		}
	}
	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    c.maxSize,
		MaxAge:     c.maxAge,
		MaxBackups: c.maxBackups,
		Compress:   c.compressed,
		LocalTime:  c.localTime,
	}
}

func makeSureDirExists(dir string) error {
	if dir != "" && dir != "." && dir != ".." {
		if _, err := os.Stat(dir); err != nil {
			if os.IsNotExist(err) {
				if err = os.MkdirAll(dir, 0755); err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}
	return nil
}

func New(cfg *elog.OutputConfig) (elog.Output, error) {
	if cfg == nil {
		return nil, errors.New("missing file output configs")
	}
	params := param.Params(cfg.Params)

	logDir := params.String(LogDir, "logs")
	baseName := params.String(BaseName, filepath.Base(os.Args[0]))

	if err := makeSureDirExists(logDir); err != nil {
		return nil, err
	}
	if strings.HasSuffix(baseName, ".log") {
		baseName = strings.TrimSuffix(baseName, ".log")
	}
	baseName = filepath.Join(logDir, baseName)

	multiFile := params.Bool(MultiFile, false)
	rc := getRollingConfig(params)

	out := &fileOutput{}

	if multiFile {
		out.lws = make([]fileWriter, 3)
		w0 := newFileLogger(fmt.Sprintf("%s.error.log", baseName), rc)
		w1 := newFileLogger(fmt.Sprintf("%s.warn.log", baseName), rc)
		w2 := newFileLogger(fmt.Sprintf("%s.info.log", baseName), rc)
		out.lws[0] = fileWriter{l: elog.ERROR, w: w0}
		out.lws[1] = fileWriter{l: elog.WARN, w: w1}
		out.lws[2] = fileWriter{l: elog.DEBUG, w: w2}
	} else {
		w0 := newFileLogger(fmt.Sprintf("%s.log", baseName), rc)
		out.lws = []fileWriter{{l: elog.DEBUG, w: w0}}
	}
	return out, nil
}
