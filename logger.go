package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/chrispassas/tformat"
)

type Logger interface {
	Debugf(string, ...interface{}) error
	Infof(string, ...interface{}) error
	Warnf(string, ...interface{}) error
	Errorf(string, ...interface{}) error
}

const (
	Restore   = "\033[0m"
	Red       = "\033[00;31m"
	Green     = "\033[00;32m"
	Yellow    = "\033[00;33m"
	Blue      = "\033[00;34m"
	Purple    = "\033[00;35m"
	Cyan      = "\033[00;36m"
	LightGrey = "\033[00;37m"
)

var writerMutex sync.Mutex
var std = newDefaultLogger()

func newDefaultLogger() *DefaultLogger {

	return &DefaultLogger{
		pid:      os.Getpid(),
		logLevel: LogLevelDebug,
		writer:   os.Stderr,
	}
}

type DefaultLogger struct {
	m                        sync.Mutex
	uuid                     string
	pid                      int
	printPid                 bool
	logLevel                 LogLevel
	writer                   io.Writer
	previousLogTime          time.Time
	printDiffPreviousLogTime bool
	printFullFilePath        bool
	disableWriterMutex       bool
	prefix                   string
	disableColor             bool
}

// TODO update so config can be used to set all values
type DefaultLoggerConfig struct {
}

type LogLevel int

const (
	LogLevelDebug LogLevel = 1
	LogLevelInfo  LogLevel = 2
	LogLevelWarn  LogLevel = 3
	LogLevelError LogLevel = 4
)

func NewDefaultLogger(config DefaultLoggerConfig) (defaultLogger *DefaultLogger) {
	defaultLogger = &DefaultLogger{
		pid:          os.Getpid(),
		logLevel:     LogLevelDebug,
		writer:       os.Stderr,
		disableColor: true,
	}

	return defaultLogger
}

func EnableColor() {
	std.EnableColor()
}

func DisableColor() {
	std.EnableColor()
}

func SetPrefix(prefix string) {
	std.SetPrefix(prefix)
}

func EnableLogDiffs() {
	std.EnableLogDiffs()
}

func DisableLogDiffs() {
	std.DisableLogDiffs()
}

func SetWriter(w io.Writer) {
	std.SetWriter(w)
}

func SetLogLevel(level LogLevel) {
	std.SetLogLevel(level)
}

func SetUUID(uuid string) {
	std.SetUUID(uuid)
}

func EnablePid() {
	std.EnablePid()
}

func DisablePid() {
	std.DisablePid()
}

func Debugf(format string, args ...interface{}) error {
	return std.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) error {
	return std.Infof(format, args...)

}

func Warnf(format string, args ...interface{}) error {
	return std.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) error {
	return std.Errorf(format, args...)
}

func (m *DefaultLogger) EnableColor() {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	m.disableColor = false
}

func (m *DefaultLogger) DisableColor() {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	m.disableColor = true
}

func (m *DefaultLogger) SetPrefix(prefix string) {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	m.prefix = prefix
}

func (m *DefaultLogger) EnableLogDiffs() {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	m.printDiffPreviousLogTime = true
}

func (m *DefaultLogger) DisableLogDiffs() {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	m.printDiffPreviousLogTime = false
}

func (m *DefaultLogger) SetWriter(w io.Writer) {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	m.writer = w
}

func (m *DefaultLogger) SetLogLevel(level LogLevel) {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	m.logLevel = level
}

func (m *DefaultLogger) SetUUID(uuid string) {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	m.uuid = uuid
}

func (m *DefaultLogger) EnablePid() {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	m.printPid = true
}

func (m *DefaultLogger) DisablePid() {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	m.printPid = false
}

func (m *DefaultLogger) Debugf(format string, args ...interface{}) error {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	if m.logLevel <= LogLevelDebug {
		return m.logf(LogLevelDebug, format, args...)
	}
	return nil
}

func (m *DefaultLogger) Infof(format string, args ...interface{}) error {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	if m.logLevel <= LogLevelInfo {
		return m.logf(LogLevelInfo, format, args...)
	}
	return nil

}

func (m *DefaultLogger) Warnf(format string, args ...interface{}) error {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	if m.logLevel <= LogLevelWarn {
		return m.logf(LogLevelWarn, format, args...)
	}
	return nil
}

func (m *DefaultLogger) Errorf(format string, args ...interface{}) error {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	if m.logLevel <= LogLevelError {
		return m.logf(LogLevelError, format, args...)
	}
	return nil
}

func (m *DefaultLogger) getExtras() (extra string) {

	if len(m.uuid) > 0 {
		extra += fmt.Sprintf(" UUID:%s", m.uuid)
	}
	if m.printPid {
		extra += fmt.Sprintf(" PID:%d", m.pid)
	}
	return extra
}

func (m *DefaultLogger) logf(logLevel LogLevel, format string, args ...interface{}) (err error) {
	var (
		msg                 string
		file                string
		line                int
		ok                  bool
		logTime             = time.Now()
		logTimeStr          = logTime.Format(fmt.Sprintf("%s %s:%s:%s%s", tformat.YYYY_MM_DD, tformat.HH24, tformat.MI, tformat.SS, tformat.Micro))
		previousLogTimeDiff time.Duration
		prefix              string
		prevLogDiffStr      string
		level               string
	)
	msg = fmt.Sprintf(format, args...)

	_, file, line, ok = runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}

	if !m.printFullFilePath {
		file = filepath.Base(file)
	}

	if m.prefix != "" {
		prefix = " (" + m.prefix + ")"
	}

	if !m.disableColor {
		switch logLevel {
		case LogLevelDebug:
			level = Purple + "DEBUG" + Restore
		case LogLevelInfo:
			level = Blue + "INFO" + Restore
		case LogLevelWarn:
			level = Yellow + "WARN" + Restore
		case LogLevelError:
			level = Red + "ERROR" + Restore
		}
	} else {
		switch logLevel {
		case LogLevelDebug:
			level = "DEBUG"
		case LogLevelInfo:
			level = "INFO"
		case LogLevelWarn:
			level = "WARN"
		case LogLevelError:
			level = "ERROR"
		}
	}

	if m.printDiffPreviousLogTime {
		if m.previousLogTime.IsZero() {
			previousLogTimeDiff = 0
		} else {
			previousLogTimeDiff = logTime.Sub(m.previousLogTime)
		}

		m.previousLogTime = logTime
		prevLogDiffStr = fmt.Sprintf(" PrevLogDiff:%s", previousLogTimeDiff)
	}
	_, err = m.writer.Write([]byte(fmt.Sprintf("%s %s:%d%s [%s] %s %s%s\n", logTimeStr, file, line, prefix, level, msg, prevLogDiffStr, m.getExtras())))

	return err
}