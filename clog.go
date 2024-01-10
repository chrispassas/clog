package clog

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type Logger interface {
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
}

// Log struct for json encoding
type Log struct {
	Time          string `json:"time,omitempty"`
	File          string `json:"file,omitempty"`
	Line          int    `json:"line,omitempty"`
	Prefix        string `json:"prefix,omitempty"`
	Level         string `json:"level,omitempty"`
	Msg           string `json:"msg,omitempty"`
	PrevLogDiffMS int64  `json:"prev_log_diff_ms,omitempty"`
	Pid           int    `json:"pid,omitempty"`
	UUID          string `json:"uuid,omitempty"`
}

type LogLevel int
type OutputFormat int

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
	LogLevelFatal LogLevel = 5

	OutputFormatStd        OutputFormat = 1
	OutputFormatJSON       OutputFormat = 2
	OutputFormatJSONIndent OutputFormat = 3

	defaultTimeFormat = "2006-01-02 15:04:05.000000"
)

// writerMutex by default this mutex is used by all instances of clog.
// This makes it safe for multiple instances of clog to write to the same io.Writer
var writerMutex sync.Mutex //nolint

// std this instance of clog can be used without creating an instance similar to the Go standard library logger
var std = newStd() //nolint

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
	outputFormat             OutputFormat
	// isStd if set to true use +1 on runtime.Caller()
	isStd bool
}

func newStd() (cLog *CLog) {
	cLog = New()
	cLog.isStd = true
	return cLog
}

// New return new instance of clog
func New() (cLog *CLog) {
	cLog = &CLog{
		pid:            os.Getpid(),
		logLevel:       LogLevelDebug,
		writer:         os.Stderr,
		disableColor:   true,
		dateTimeFormat: defaultTimeFormat,
		outputFormat:   OutputFormatStd,
	}

	return cLog
}

// SetOutputFormat change output format for standard clog instance
func SetOutputFormat(format OutputFormat) {
	std.SetOutputFormat(format)
}

// SetTimeFormat change time format for standard clog instance
func SetTimeFormat(format string) {
	std.SetTimeFormat(format)
}

// EnableColor turn on color output for standard clog instance
func EnableColor() {
	std.EnableColor()
}

// DisableColor turn off color output for standard clog instance
func DisableColor() {
	std.EnableColor()
}

// SetPrefix set a prefix for output for standard clog instance
func SetPrefix(prefix string) {
	std.SetPrefix(prefix)
}

// EnableLogDiffs turn on time diff since last log line for standard clog instance
// This can makes it easy to see how much time has passed since the last log line.
func EnableLogDiffs() {
	std.EnableLogDiffs()
}

// DisableLogDiffs turn off log time diffing for standard clog instance
func DisableLogDiffs() {
	std.DisableLogDiffs()
}

// SetWriter set io.Writer for output for standard clog instance
func SetWriter(w io.Writer) {
	std.SetWriter(w)
}

// SetLogLevel set log level (clog.LogLevelDebug, clog.LogLevelInfo, clog.LogLevelWarn, clog.LogLevelError) for standard clog instance
func SetLogLevel(level LogLevel) {
	std.SetLogLevel(level)
}

// SetUUID set UUID to print at the end of each line for standard clog instance
func SetUUID(uuid string) {
	std.SetUUID(uuid)
}

// EnablePid turn on printing pid to end of each line for standard clog instance
// This can be useful if multiple processes could be writing to the same log file (service reload)
func EnablePid() {
	std.EnablePid()
}

// DisablePid turn off printing pid to end of file for standard clog instance
func DisablePid() {
	std.DisablePid()
}

// Debugf prints DEBUG level for standard clog instance
func Debugf(format string, args ...interface{}) {
	std.Debugf(format, args...)
}

// Infof prints INFO level for standard clog instance
func Infof(format string, args ...interface{}) {
	std.Infof(format, args...)

}

// Warnf prints WARN level for standard clog instance
func Warnf(format string, args ...interface{}) {
	std.Warnf(format, args...)
}

// Errorf prints ERROR level for standard clog instance
func Errorf(format string, args ...interface{}) {
	std.Errorf(format, args...)
}

// Errorf prints FATAL level for standard clog instance
func Fatalf(format string, args ...interface{}) {
	std.Fatalf(format, args...)
}

// SetOutputFormat change output format
func (m *CLog) SetOutputFormat(format OutputFormat) {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	m.outputFormat = format
}

// SetTimeFormat change time format
func (m *CLog) SetTimeFormat(format string) {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	m.dateTimeFormat = format
}

// EnableColor turn on color output
func (m *CLog) EnableColor() {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	m.disableColor = false
}

// DisableColor turn off color output
func (m *CLog) DisableColor() {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	m.disableColor = true
}

// SetPrefix set a prefix
func (m *CLog) SetPrefix(prefix string) {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	m.prefix = prefix
}

// EnableLogDiffs turn on time diff since last log line
// This can makes it easy to see how much time has passed since the last log line.
func (m *CLog) EnableLogDiffs() {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	m.printDiffPreviousLogTime = true
}

// DisableLogDiffs turn off log time diffing
func (m *CLog) DisableLogDiffs() {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	m.printDiffPreviousLogTime = false
}

// SetWriter set io.Writer for output
func (m *CLog) SetWriter(w io.Writer) {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	m.writer = w
}

// SetLogLevel set log level (clog.LogLevelDebug, clog.LogLevelInfo, clog.LogLevelWarn, clog.LogLevelError)
func (m *CLog) SetLogLevel(level LogLevel) {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	m.logLevel = level
}

// SetUUID set UUID to print at the end of each line
func (m *CLog) SetUUID(uuid string) {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	m.uuid = uuid
}

// EnablePid turn on printing pid to end of each line
// This can be useful if multiple processes could be writing to the same log file (service reload)
func (m *CLog) EnablePid() {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	m.printPid = true
}

// DisablePid turn off printing pid to end of file
func (m *CLog) DisablePid() {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	m.printPid = false
}

// Debugf prints DEBUG level
func (m *CLog) Debugf(format string, args ...interface{}) {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	if m.logLevel <= LogLevelDebug {
		m.logf(LogLevelDebug, format, args...)
	}
}

// Infof prints INFO level
func (m *CLog) Infof(format string, args ...interface{}) {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	if m.logLevel <= LogLevelInfo {
		m.logf(LogLevelInfo, format, args...)
	}
}

// Warnf prints WARN level
func (m *CLog) Warnf(format string, args ...interface{}) {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	if m.logLevel <= LogLevelWarn {
		m.logf(LogLevelWarn, format, args...)
	}
}

// Errorf prints ERROR level
func (m *CLog) Errorf(format string, args ...interface{}) {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	if m.logLevel <= LogLevelError {
		m.logf(LogLevelError, format, args...)
	}
}

// Fatalf prints FATAL level then os.Exit(1)
func (m *CLog) Fatalf(format string, args ...interface{}) {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	if m.logLevel <= LogLevelFatal {
		m.logf(LogLevelFatal, format, args...)
	}
	os.Exit(1) //nolint
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

	if m.isStd {
		_, file, line, ok = runtime.Caller(3)
	} else {
		_, file, line, ok = runtime.Caller(2)
	}

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

	if m.printDiffPreviousLogTime {
		if m.previousLogTime.IsZero() {
			previousLogTimeDiff = 0
		} else {
			previousLogTimeDiff = logTime.Sub(m.previousLogTime)
		}

		m.previousLogTime = logTime
		prevLogDiffStr = fmt.Sprintf(" PrevLogDiff:%s", previousLogTimeDiff)
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
	case LogLevelFatal:
		level = "FATAL"
	}

	switch m.outputFormat {
	case OutputFormatStd:

		if !m.disableColor {
			switch logLevel {
			case LogLevelDebug:
				level = purple + level + restore
			case LogLevelInfo:
				level = blue + level + restore
			case LogLevelWarn:
				level = yellow + level + restore
			case LogLevelError, LogLevelFatal:
				level = red + level + restore
			}
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
		break
	case OutputFormatJSON, OutputFormatJSONIndent:

		log := Log{
			Time:          logTimeStr,
			File:          file,
			Line:          line,
			Prefix:        prefix,
			Level:         level,
			Msg:           msg,
			PrevLogDiffMS: previousLogTimeDiff.Milliseconds(),
		}
		if len(m.uuid) > 0 {
			log.UUID = m.uuid
		}
		if m.printPid {
			log.Pid = m.pid
		}

		var jsonBytes []byte
		switch m.outputFormat {
		case OutputFormatJSON:
			if jsonBytes, err = json.Marshal(&log); err != nil {
				err = fmt.Errorf("json.Marshal() error:%w", err)
				return err
			}
		case OutputFormatJSONIndent:
			if jsonBytes, err = json.MarshalIndent(&log, "", "\t"); err != nil {
				err = fmt.Errorf("json.Marshal() error:%w", err)
				return err
			}
		}

		if n, err = m.writer.Write(append(jsonBytes, '\n')); err != nil {
			err = fmt.Errorf("m.writer.Write() n:%d error:%w", n, err)
		}
		break
	default:
		err = fmt.Errorf("Unsupported OutputFormat:%d", m.outputFormat)
	}

	return err
}
