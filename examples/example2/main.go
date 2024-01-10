package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chrispassas/clog"
)

func main() {
	logger := clog.New()
	logger.SetPrefix("main")
	logger.EnableLogDiffs()
	logger.EnablePid()
	logger.Infof("Main started")

	// Define the endpoint handler function
	handleTest := func(w http.ResponseWriter, r *http.Request) {
		requestLogger := clog.New()
		requestLogger.EnableLogDiffs()
		requestLogger.EnablePid()
		requestLogger.SetPrefix("handleTest")
		requestLogger.SetUUID("0d01be9f-f965-4398-a046-1e83322cb243")
		requestLogger.Debugf("request recieved")
		time.Sleep(time.Second * 2)
		fmt.Fprintln(w, "Example HTTP Response")
		requestLogger.Debugf("request complete")
	}

	// Register the handler function for the "/test" endpoint
	http.HandleFunc("/test", handleTest)

	// Start the web server on port 8080
	logger.Infof("Server is listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		logger.Errorf("http.ListenAndServe() error:%v", err)
	}
}
