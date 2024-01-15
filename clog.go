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
	// Time when the log method was called Debugf, Infof, Warnf, Errorf
	Time string `json:"time,omitempty"`
	// File that called log method
	File string `json:"file,omitempty"`
	// Line number where log method was called
	Line int `json:"line,omitempty"`
	// Prefix allows adding a label to different instances of the logger
	Prefix string `json:"prefix,omitempty"`
	// Level debug, info, warn, error
	Level string `json:"level,omitempty"`
	// Msg log message
	Msg string `json:"msg,omitempty"`
	// Diff time passes since previous log method was called
	Diff int64 `json:"diff,omitempty"`
	// Pid process id
	Pid int `json:"pid,omitempty"`
	// UUID optional uuid for tracing and other uses
	UUID string `json:"uuid,omitempty"`
}

type LogLevel int
type OutputFormat int
type PrintSource int

const (
	restore   = "\033[0m"
	red       = "\033[00;31m"
	green     = "\033[00;32m"
	yellow    = "\033[00;33m"
	blue      = "\033[00;34m"
	purple    = "\033[00;35m"
	cyan      = "\033[00;36m"
	lightGrey = "\033[00;37m"

	PrintSourceDisable  PrintSource = 1
	PrintSourceFile     PrintSource = 2
	PrintSourceFullPath PrintSource = 3

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
	printSource              PrintSource
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
		printSource:    PrintSourceFile,
		dateTimeFormat: defaultTimeFormat,
		outputFormat:   OutputFormatStd,
	}

	return cLog
}

// Copy returns a copy of the current clog, changes to this copy will not affect the original clog instance
func Copy() *CLog {
	newClog := std.Copy()
	return newClog
}

// SetPrintSource change between for standard clog instance
// PrintSourceDisable - don't print source file and line number
// PrintSourceFile - print just source file and line number
// PrintSourceFullPath - print full path with source file and line number
func SetPrintSource(printSource PrintSource) *CLog {
	std.SetPrintSource(printSource)
	return std
}

// SetOutputFormat change output format for standard clog instance
func SetOutputFormat(format OutputFormat) *CLog {
	std.SetOutputFormat(format)
	return std
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

// EnableDiffs turn on time diff since last log line for standard clog instance
// This can makes it easy to see how much time has passed since the last log line.
func EnableDiffs() *CLog {
	std.EnableDiffs()
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

// Copy returns a copy of the current clog, changes to this copy will not affect the original clog instance
func (m *CLog) Copy() *CLog {
	newCLog := &CLog{
		uuid:                     m.uuid,
		pid:                      m.pid,
		printPid:                 m.printPid,
		logLevel:                 m.logLevel,
		writer:                   m.writer,
		printDiffPreviousLogTime: m.printDiffPreviousLogTime,
		printSource:              m.printSource,
		disableWriterMutex:       m.disableWriterMutex,
		prefix:                   m.prefix,
		disableColor:             m.disableColor,
		dateTimeFormat:           m.dateTimeFormat,
		outputFormat:             m.outputFormat,
	}

	return newCLog
}

// SetPrintSource change between
// PrintSourceDisable - don't print source file and line number
// PrintSourceFile - print just source file and line number
// PrintSourceFullPath - print full path with source file and line number
func (m *CLog) SetPrintSource(printSource PrintSource) *CLog {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	m.printSource = printSource
	return m
}

// SetOutputFormat change output format
func (m *CLog) SetOutputFormat(format OutputFormat) *CLog {
	if !m.disableWriterMutex {
		writerMutex.Lock()
		defer writerMutex.Unlock()
	} else {
		m.m.Lock()
		defer m.m.Unlock()
	}
	m.outputFormat = format
	return m
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

// EnableDiffs turn on time diff since last log line
// This can makes it easy to see how much time has passed since the last log line.
func (m *CLog) EnableDiffs() *CLog {
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
		fileStr             string
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

	switch m.printSource {
	case PrintSourceDisable:
		fileStr = ""
		file = ""
	case PrintSourceFile:
		file = filepath.Base(file)
		fileStr = fmt.Sprintf("%s:%d", file, line)
	case PrintSourceFullPath:
		fileStr = fmt.Sprintf("%s:%d", file, line)
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
		prevLogDiffStr = fmt.Sprintf(" DIFF:%s", previousLogTimeDiff)
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
			[]byte(fmt.Sprintf("%s %s%s [%s] %s%s%s\n",
				logTimeStr,
				fileStr,
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
			Time:   logTimeStr,
			File:   file,
			Line:   line,
			Prefix: prefix,
			Level:  level,
			Msg:    msg,
			Diff:   previousLogTimeDiff.Milliseconds(),
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
