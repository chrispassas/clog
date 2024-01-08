package main

import (
	"fmt"
	"log"

	"github.com/chrispassas/logger"
)

func main() {
	log.Printf("test")
	log := logger.NewDefaultLogger(logger.DefaultLoggerConfig{})
	// log.EnableLogDiffs()
	// log.SetPrefix("main")
	var err error
	if err = log.Debugf("test"); err != nil {
		fmt.Printf("error:%v", err)
	}
	log.Infof("test2")
	log.Warnf("test3")
	log.Errorf("test4")
}
