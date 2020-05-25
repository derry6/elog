package bench

import (
	"fmt"

	"github.com/derry6/elog/output/file"
	"github.com/stretchr/testify/assert"

	"io"
	"io/ioutil"
	"testing"

	"github.com/derry6/elog"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	encodingText = "text"
	encodingJSON = "json"
)

var (
	szG      = 1024
	_maxSize = 4 * szG
	_maxAge  = 1
	_backups = 0
)

func testConsoleConfig() *elog.OutputConfig {
	cfg := elog.DefaultConsoleConfig()
	cfg.Enabled = true
	cfg.Params["writer"] = "discard"
	return cfg
}

func buildFileName(logName string, encoding string) string {
	return fmt.Sprintf("logs/%s_%s.log", logName, encoding)
}

func rollingPolicy() map[string]interface{} {
	values := make(map[string]interface{})
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
	logger := newZapLogger(encodingText, zapcore.DebugLevel, w)
	benchmark(b, logger.Sugar())
}
func BenchmarkZapJson(b *testing.B) {
	w := zapcore.Lock(zapcore.AddSync(ioutil.Discard))
	logger := newZapLogger(encodingJSON, zap.DebugLevel, w)
	benchmark(b, logger.Sugar())
}
func BenchmarkZapTextFile(b *testing.B) {
	w := toLogFile(buildFileName("zap", encodingText))
	logger := newZapLogger(encodingText, zap.DebugLevel, zapcore.AddSync(w))
	sugar := logger.Sugar()
	benchmark(b, sugar)
}
func BenchmarkZapJsonFile(b *testing.B) {
	w := toLogFile(buildFileName("zap", encodingJSON))
	logger := newZapLogger(encodingJSON, zap.DebugLevel, zapcore.AddSync(w))
	sugar := logger.Sugar()
	benchmark(b, sugar)
}

func benchmarkLoggerFile(encoding string, b *testing.B) {
	logger := elog.New(nil)
	err := logger.AddOutput(file.Name, &elog.OutputConfig{
		Enabled: true,
		Params:  rollingPolicy(),
	})
	assert.NoError(b, err)
	benchmark(b, logger)
	_ = logger.Sync()
}

func BenchmarkELogText(b *testing.B) {
	cfg := testConsoleConfig()
	cfg.Encoding = elog.Text
	logger := elog.New(nil)
	assert.NoError(b, logger.AddOutput(elog.Console, cfg))
	benchmark(b, logger)
}
func BenchmarkELogJson(b *testing.B) {
	cfg := testConsoleConfig()
	cfg.Encoding = elog.JSON
	logger := elog.New(nil)
	assert.NoError(b, logger.AddOutput(elog.Console, cfg))
	benchmark(b, logger)
}

func BenchmarkELogTextFile(b *testing.B) {
	benchmarkLoggerFile(encodingText, b)
}
func BenchmarkELogJsonFile(b *testing.B) {
	benchmarkLoggerFile(encodingJSON, b)
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
	logger := newLogrusLogger(encodingText, ioutil.Discard)
	logger.SetReportCaller(false)
	benchmark(b, &logrusLogger{logger})
}
func BenchmarkLogrusJson(b *testing.B) {
	logger := newLogrusLogger(encodingJSON, ioutil.Discard)
	logger.SetReportCaller(false)
	benchmark(b, &logrusLogger{logger})
}
func BenchmarkLogrusTextFile(b *testing.B) {
	w := toLogFile(buildFileName("logrus", encodingText))
	logger := newLogrusLogger(encodingText, w)
	logger.SetReportCaller(false)
	logger.SetNoLock()
	benchmark(b, &logrusLogger{logger})
}
func BenchmarkLogrusJsonFile(b *testing.B) {
	w := toLogFile(buildFileName("logrus", encodingJSON))
	logger := newLogrusLogger(encodingJSON, w)
	logger.SetReportCaller(false)
	logger.SetNoLock()
	benchmark(b, &logrusLogger{logger})
}
