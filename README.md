[![GoDoc](https://godoc.org/github.com/chrispassas/clog?status.svg)](https://godoc.org/github.com/chrispassas/clog)

# clog
This is a simple logger that is easy to understand. Only stdlib libraries are used.

## Installation

`go get -u github.com/chrispassas/clog`

## Quick Start
```go
	clog.Debugf("This is a debug message")
	clog.Infof("This is a info message")
	clog.Warnf("This is a warn message")
	clog.Errorf("This is a error message")

	// Or create instance of logger you can customize

	log := clog.New()
	log.SetPrefix("main")
	log.EnablePid()
	log.Infof("Main started")
```

See the [documentation][doc]

## Example 1
Global clog

In this example you can use clog similar to how you might use the standard package log.

```go
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

	clog.EnableDiffs()
	clog.Debugf("printing log with previous line diff turned on")
	time.Sleep(time.Second)
	clog.Debugf("printing log with previous line diff turned on again")

	clog.EnablePid() // Add process ID to log line
	clog.SetPrefix("main")
	clog.Debugf("This has a prefix and pid")
}

```
### Output
```bash
2024-01-09 17:55:42.270909 main.go:14 [DEBUG] This is a debug message
2024-01-09 17:55:42.271504 main.go:15 [INFO] This is a info message
2024-01-09 17:55:42.271514 main.go:16 [WARN] This is a warn message
2024-01-09 17:55:42.271520 main.go:17 [ERROR] This is a error message
2024-01-09 17:55:42.271526 main.go:20 [DEBUG] printing log with previous line diff turned on DIFF:0s
2024-01-09 17:55:43.272648 main.go:22 [DEBUG] printing log with previous line diff turned on again DIFF:1.001120708s
2024-01-09 17:55:43.272959 main.go:26 (main) [DEBUG] This has a prefix and pid DIFF:311.667µs PID:12462
```

## Example 2
Instance of clog

This shows creating multiple instances of clog so the output can be customized depending on the codes needs.
For example different goroutines might want a different log prefix or UUID per-request.

Each instance can also set its own io.Writer (output file) as needed.

```go
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
	logger.EnableDiffs()
	logger.EnablePid()
	logger.Infof("Main started")

	// Define the endpoint handler function
	handleTest := func(w http.ResponseWriter, r *http.Request) {
		requestLogger := clog.New()
		requestLogger.EnableDiffs()
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

```

### Output
```bash
2024-01-09 18:03:33.905305 main.go:14 (main) [INFO] Main started DIFF:0s PID:13308
2024-01-09 18:03:33.905683 main.go:29 (main) [INFO] Server is listening on :8080 DIFF:378.208µs PID:13308
2024-01-09 18:03:46.845565 main.go:19 (handleTest) [DEBUG] request recieved DIFF:0s uuid:0d01be9f-f965-4398-a046-1e83322cb243 PID:13308
2024-01-09 18:03:48.848360 main.go:22 (handleTest) [DEBUG] request complete DIFF:2.0028025s UUID:0d01be9f-f965-4398-a046-1e83322cb243 PID:13308
```

[doc]: https://godoc.org/github.com/chrispassas/clog