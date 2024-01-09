package clog

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type Logger interface {
	Debugf(string, ...interface{}) error
	Infof(string, ...interface{}) error
	Warnf(string, ...interface{}) error
	Errorf(string, ...interface{}) error
}

type LogLevel int

const (
	restore   = "\033[0m"
	red       = "\033[00;31m"
	green     = "\033[00;32m"
	yellow    = "\033[00;33m"
	blue      = "\033[00;34m"
	purple    = "\033[00;35m"
	cyan      = "\033[00;36m"
	lightGrey = "\033[00;37m"

	LogLevelDebug LogLevel = 1
	LogLevelInfo  LogLevel = 2
	LogLevelWarn  LogLevel = 3
	LogLevelError LogLevel = 4

	defaultTimeFormat = "2006-01-02 15:04:05.000000"
)

// writerMutex by default this mutex is used by all instances of clog.
// This makes it safe for multiple instances of clog to write to the same io.Writer
var writerMutex sync.Mutex //nolint

// std this instance of clog can be used without creating an instance similar to the Go standard library logger
var std = New() //nolint

// CLog contains all settings and values to support clog usage
type CLog struct {
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
	dateTimeFormat           string
}

// New return new instance of clog
func New() (cLog *CLog) {
	cLog = &CLog{
		pid:            os.Getpid(),
		logLevel:       LogLevelDebug,
		writer:         os.Stderr,
		disableColor:   true,
		dateTimeFormat: defaultTimeFormat,
	}

	return cLog
}

// SetTimeFormat change time format for standard clog instance
func SetTimeFormat(format string) *CLog {
	std.SetTimeFormat(format)
	return std
}

// EnableColor turn on color output for standard clog instance
func EnableColor() *CLog {
	std.EnableColor()
	return std
}

// DisableColor turn off color output for standard clog instance
func DisableColor() *CLog {
	std.EnableColor()
	return std
}

// SetPrefix set a prefix for output for standard clog instance
func SetPrefix(prefix string) *CLog {
	std.SetPrefix(prefix)
	return std
}

// EnableLogDiffs turn on time diff since last log line for standard clog instance
// This can makes it easy to see how much time has passed since the last log line.
func EnableLogDiffs() *CLog {
	std.EnableLogDiffs()
	return std
}

// DisableLogDiffs turn off log time diffing for standard clog instance
func DisableLogDiffs() *CLog {
	std.DisableLogDiffs()
	return std
}

// SetWriter set io.Writer for output for standard clog instance
func SetWriter(w io.Writer) *CLog {
	std.SetWriter(w)
	return std
}

// SetLogLevel set log level (clog.LogLevelDebug, clog.LogLevelInfo, clog.LogLevelWarn, clog.LogLevelError) for standard clog instance
func SetLogLevel(level LogLevel) *CLog {
	std.SetLogLevel(level)
	return std
}

// SetUUID set UUID to print at the end of each line for standard clog instance
func SetUUID(uuid string) *CLog {
	std.SetUUID(uuid)
	return std
}

// EnablePid turn on printing pid to end of each line for standard clog instance
// This can be useful if multiple processes could be writing to the same log file (service reload)
func EnablePid() *CLog {
	std.EnablePid()
	return std
}

// DisablePid turn off printing pid to end of file for standard clog instance
func DisablePid() *CLog {
	std.DisablePid()
	return std
}

// Debugf prints DEBUG level for standard clog instance
func Debugf(format string, args ...interface{}) error {
	return std.Debugf(format, args...)
}

// Infof prints INFO level for standard clog instance
func Infof(format string, args ...interface{}) error {
	return std.Infof(format, args...)

}

// Warnf prints WARN level for standard clog instance
func Warnf(format string, args ...interface{}) error {
	return std.Warnf(format, args...)
}

// Errorf prints ERROR level for standard clog instance
func Errorf(format string, args ...interface{}) error {
	return std.Errorf(format, args...)
}

// SetTimeFormat change time format
func (m *CLog) SetTimeFormat(format string) *CLog {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	m.dateTimeFormat = format
	return m
}

// EnableColor turn on color output
func (m *CLog) EnableColor() *CLog {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	m.disableColor = false
	return m
}

// DisableColor turn off color output
func (m *CLog) DisableColor() *CLog {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	m.disableColor = true
	return m
}

// SetPrefix set a prefix
func (m *CLog) SetPrefix(prefix string) *CLog {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	m.prefix = prefix
	return m
}

// EnableLogDiffs turn on time diff since last log line
// This can makes it easy to see how much time has passed since the last log line.
func (m *CLog) EnableLogDiffs() *CLog {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	m.printDiffPreviousLogTime = true
	return m
}

// DisableLogDiffs turn off log time diffing
func (m *CLog) DisableLogDiffs() *CLog {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	m.printDiffPreviousLogTime = false
	return m
}

// SetWriter set io.Writer for output
func (m *CLog) SetWriter(w io.Writer) *CLog {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	m.writer = w
	return m
}

// SetLogLevel set log level (clog.LogLevelDebug, clog.LogLevelInfo, clog.LogLevelWarn, clog.LogLevelError)
func (m *CLog) SetLogLevel(level LogLevel) *CLog {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	m.logLevel = level
	return m
}

// SetUUID set UUID to print at the end of each line
func (m *CLog) SetUUID(uuid string) *CLog {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	m.uuid = uuid
	return m
}

// EnablePid turn on printing pid to end of each line
// This can be useful if multiple processes could be writing to the same log file (service reload)
func (m *CLog) EnablePid() *CLog {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	m.printPid = true
	return m
}

// DisablePid turn off printing pid to end of file
func (m *CLog) DisablePid() *CLog {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	m.printPid = false
	return m
}

// Debugf prints DEBUG level
func (m *CLog) Debugf(format string, args ...interface{}) error {
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

// Infof prints INFO level
func (m *CLog) Infof(format string, args ...interface{}) error {
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

// Warnf prints WARN level
func (m *CLog) Warnf(format string, args ...interface{}) error {
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

// Errorf prints ERROR level
func (m *CLog) Errorf(format string, args ...interface{}) error {
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

func (m *CLog) getExtras() (extra string) {

	if len(m.uuid) > 0 {
		extra += fmt.Sprintf(" UUID:%s", m.uuid)
	}
	if m.printPid {
		extra += fmt.Sprintf(" PID:%d", m.pid)
	}
	return extra
}

func (m *CLog) logf(logLevel LogLevel, format string, args ...interface{}) (err error) {
	var (
		msg                 string
		file                string
		line                int
		ok                  bool
		logTime             = time.Now()
		logTimeStr          = logTime.Format(m.dateTimeFormat)
		previousLogTimeDiff time.Duration
		prefix              string
		prevLogDiffStr      string
		level               string
		n                   int
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

	if !m.disableColor {
		switch logLevel {
		case LogLevelDebug:
			level = purple + level + restore
		case LogLevelInfo:
			level = blue + level + restore
		case LogLevelWarn:
			level = yellow + level + restore
		case LogLevelError:
			level = red + level + restore
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

	if n, err = m.writer.Write(
		[]byte(fmt.Sprintf("%s %s:%d%s [%s] %s%s%s\n",
			logTimeStr,
			file,
			line,
			prefix,
			level,
			msg,
			prevLogDiffStr,
			m.getExtras(),
		))); err != nil {
		err = fmt.Errorf("m.writer.Write() n:%d error:%w", n, err)
	}

	return err
}
