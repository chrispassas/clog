package main

import (
	"time"

	"github.com/chrispassas/clog"
)

func main() {
	clog.Debugf("This is a debug message")
	clog.Infof("This is a info message")
	clog.Warnf("This is a warn message")
	clog.Errorf("This is a error message")

	clog.EnableLogDiffs()
	clog.Debugf("printing log with previous line diff turned on")
	time.Sleep(time.Second)
	clog.Debugf("printing log with previous line diff turned on again")

	clog.EnablePid() // Add process ID to log line
	clog.SetPrefix("main")
	clog.Debugf("This has a prefix and pid")
}
