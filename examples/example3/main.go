package main

import (
	"github.com/chrispassas/clog"
)

func main() {
	clog.Debugf("This is a debug message")
	clog.Infof("This is a info message")
	clog.Warnf("This is a warn message")
	clog.Errorf("This is a error message")

	clog.Fatalf("emergancy exit")
}
