# clog
This is a simple logger that is easy to understand and use. Only stdlib libraries are used.



## Example
Global clog

```go
package main

import (
	"log"
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
```
### Output
```bash
2024-01-09 17:55:42.270909 main.go:14 [DEBUG] This is a debug message
2024-01-09 17:55:42.271504 main.go:15 [INFO] This is a info message
2024-01-09 17:55:42.271514 main.go:16 [WARN] This is a warn message
2024-01-09 17:55:42.271520 main.go:17 [ERROR] This is a error message
2024-01-09 17:55:42.271526 main.go:20 [DEBUG] printing log with previous line diff turned on PrevLogDiff:0s
2024-01-09 17:55:43.272648 main.go:22 [DEBUG] printing log with previous line diff turned on again PrevLogDiff:1.001120708s
2024-01-09 17:55:43.272959 main.go:26 (main) [DEBUG] This has a prefix and pid PrevLogDiff:311.667µs PID:12462
```

## Example
Instance of clog

```go
package main

import (
	"log"

	"github.com/chrispassas/clog"
)

func main() {
    // Example here
    log := clog.New()

}
```
