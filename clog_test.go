package clog

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestSetPrefix(t *testing.T) {
	var buffer bytes.Buffer

	logger := New()
	logger.SetWriter(&buffer)
	logger.SetPrefix("test")
	logger.Debugf("test message")

	if !strings.Contains(buffer.String(), "(test)") {
		t.Errorf("failed to print expected prefix")
	}
}

func TestSetUUID(t *testing.T) {
	var buffer bytes.Buffer
	var uuid = "0d01be9f-f965-4398-a046-1e83322cb243"
	logger := New()
	logger.SetWriter(&buffer)
	logger.SetUUID(uuid)
	logger.Debugf("test message")

	if !strings.Contains(buffer.String(), uuid) {
		t.Errorf("failed to print expected uuid")
	}
}

func TestLogging(t *testing.T) {
	logger := New()
	logger.SetOutputFormat(OutputFormatJSON)
	var buffer bytes.Buffer
	logger.SetWriter(&buffer)

	logger.Debugf("debug")
	logger.Infof("info")
	logger.Warnf("warn")
	logger.Errorf("error")
	lines := bytes.Split(buffer.Bytes(), []byte("\n"))
	for x, line := range lines {
		if len(line) == 0 {
			continue
		}
		t.Logf("line:%s", string(line))
		var log Log
		var err error
		if err = json.Unmarshal(line, &log); err != nil {
			t.Errorf("json.Unmarshal() error:%v", err)
		}
		switch x {
		case 0:
			if log.Msg != "debug" {
				t.Errorf("log msg is:%s expected 'debug'", log.Msg)
			}
			if log.Level != "DEBUG" {
				t.Errorf("log level is:%s expected 'DEBUG'", log.Level)
			}
		case 1:
			if log.Msg != "info" {
				t.Errorf("log msg is:%s expected 'info'", log.Msg)
			}
			if log.Level != "INFO" {
				t.Errorf("log level is:%s expected 'INFO'", log.Level)
			}
		case 2:
			if log.Msg != "warn" {
				t.Errorf("log msg is:%s expected 'warn'", log.Msg)
			}
			if log.Level != "WARN" {
				t.Errorf("log level is:%s expected 'WARN'", log.Level)
			}
		case 3:
			if log.Msg != "error" {
				t.Errorf("log msg is:%s expected 'error'", log.Msg)
			}
			if log.Level != "ERROR" {
				t.Errorf("log level is:%s expected 'ERROR'", log.Level)
			}
		}

	}

}
