package main

import (
	"fmt"
	"log"

	"github.com/chrispassas/clog"
)

func main() {
	log.Printf("test")
	log := clog.New().EnableColor()
	log.EnableLogDiffs()
	log.SetPrefix("main")
	log.EnablePid()
	log.SetUUID("0000")
	// log.EnableColor()
	var err error
	if err = log.Debugf("test"); err != nil {
		fmt.Printf("error:%v", err)
	}
	log.Infof("test2")
	log.Warnf("test3")
	log.Errorf("test4")

	clog.EnableColor().EnableLogDiffs().EnablePid().SetLogLevel(clog.LogLevelDebug)
	// logger.EnableLogDiffs()
	// logger.SetPrefix("foo")
	// logger.EnableColor()
	clog.Debugf("test")

}
