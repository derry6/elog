package bench

import (
	"fmt"

	"github.com/derry6/elog/internal/param"
	"github.com/derry6/elog/output/file"

	"github.com/derry6/elog/encoder"
	"github.com/derry6/elog/ezap"
	"github.com/derry6/elog/output"

	"io"
	"io/ioutil"
	"testing"

	"github.com/derry6/elog"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	szG      = 1024
	_maxSize = 4 * szG
	_maxAge  = 1
	_backups = 0
)

func buildFileName(logName string, encoding string) string {
	return fmt.Sprintf("logs/%s_%s.log", logName, encoding)
}

func rollingPolicy() param.Params {
	values := make(param.Params)
	values[file.RollingEnabled] = true
	values[file.MaxSize] = _maxSize
	values[file.MaxDays] = _maxAge
	values[file.Compressed] = false
	values[file.NameUseUTC] = false
	values[file.MaxBackups] = _backups
	return values
}

func toLogFile(name string) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   name,
		MaxSize:    _maxSize,
		MaxAge:     _maxAge,
		MaxBackups: _backups,
		LocalTime:  true,
		Compress:   false,
	}
}

type testLogger interface {
	Info(msg ...interface{})
	Infow(msg string, kvs ...interface{})
}

func benchmark(b *testing.B, logger testLogger) {
	b.Run("NoFields", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("10Fields", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Infow("test logger with 10 fields", fakeSugarFields()...)
			}
		})
	})
}
func BenchmarkZapText(b *testing.B) {
	w := zapcore.Lock(zapcore.AddSync(ioutil.Discard))
	logger := newZapLogger(encoder.Text, zapcore.DebugLevel, w)
	benchmark(b, logger.Sugar())
}
func BenchmarkZapJson(b *testing.B) {
	w := zapcore.Lock(zapcore.AddSync(ioutil.Discard))
	logger := newZapLogger(encoder.Json, zap.DebugLevel, w)
	benchmark(b, logger.Sugar())
}
func BenchmarkZapTextFile(b *testing.B) {
	w := toLogFile(buildFileName("zap", encoder.Text))
	logger := newZapLogger(encoder.Text, zap.DebugLevel, zapcore.AddSync(w))
	sugar := logger.Sugar()
	benchmark(b, sugar)
}
func BenchmarkZapJsonFile(b *testing.B) {
	w := toLogFile(buildFileName("zap", encoder.Json))
	logger := newZapLogger(encoder.Json, zap.DebugLevel, zapcore.AddSync(w))
	sugar := logger.Sugar()
	benchmark(b, sugar)
}

func benchmarkLoggerFile(encoding string, b *testing.B) {
	fileOpts := []output.Option{
		output.WithParams(rollingPolicy()),
	}
	elog.Use(
		ezap.New(
			elog.WithConsoleDisabled(),
			elog.WithDefaultEncoding(encoding),
			elog.WithFile(buildFileName("elog", encoding), fileOpts...),
		),
	)
	benchmark(b, elog.Get())
	_ = elog.Sync()
}

func BenchmarkELogText(b *testing.B) {
	elog.Use(
		ezap.New(
			elog.WithDefaultEncoding(encoder.Text),
			elog.WithConsole("discard")),
	)
	benchmark(b, elog.Get())
}
func BenchmarkELogJson(b *testing.B) {
	elog.Use(
		ezap.New(
			elog.WithDefaultEncoding(encoder.Json),
			elog.WithConsole("discard")),
	)
	benchmark(b, elog.Get())
}

func BenchmarkELogTextFile(b *testing.B) {
	benchmarkLoggerFile(encoder.Text, b)
}
func BenchmarkELogJsonFile(b *testing.B) {
	benchmarkLoggerFile(encoder.Json, b)
}

var (
	logrusFields = logrus.Fields{
		"int":     _tenInts[0],
		"ints":    _tenInts,
		"string":  _tenStrings[0],
		"strings": _tenStrings,
		"time":    _tenTimes[0],
		"times":   _tenTimes,
		"user1":   _oneUser,
		"user2":   _oneUser,
		"users":   _tenUsers,
		"error":   errExample,
	}
)

type logrusLogger struct {
	logger *logrus.Logger
}

func (l *logrusLogger) Info(v ...interface{}) {
	l.logger.Info(v...)
}
func (l *logrusLogger) Infow(msg string, v ...interface{}) {
	l.logger.WithFields(logrusFields).Info(msg)
}
func newLogrusLogger(encoding string, w io.Writer) *logrus.Logger {
	logger := logrus.New()
	logger.SetReportCaller(false)
	var formatter logrus.Formatter
	switch encoding {
	case "json":
		formatter = &logrus.JSONFormatter{}
	default:
		formatter = &logrus.TextFormatter{}
	}
	logger.SetFormatter(formatter)
	logger.SetOutput(w)
	return logger
}
func BenchmarkLogrusText(b *testing.B) {
	logger := newLogrusLogger(encoder.Text, ioutil.Discard)
	logger.SetReportCaller(false)
	benchmark(b, &logrusLogger{logger})
}
func BenchmarkLogrusJson(b *testing.B) {
	logger := newLogrusLogger(encoder.Json, ioutil.Discard)
	logger.SetReportCaller(false)
	benchmark(b, &logrusLogger{logger})
}
func BenchmarkLogrusTextFile(b *testing.B) {
	w := toLogFile(buildFileName("logrus", encoder.Text))
	logger := newLogrusLogger(encoder.Text, w)
	logger.SetReportCaller(false)
	logger.SetNoLock()
	benchmark(b, &logrusLogger{logger})
}
func BenchmarkLogrusJsonFile(b *testing.B) {
	w := toLogFile(buildFileName("logrus", encoder.Json))
	logger := newLogrusLogger(encoder.Json, w)
	logger.SetReportCaller(false)
	logger.SetNoLock()
	benchmark(b, &logrusLogger{logger})
}
