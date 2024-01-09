# clog
This is a simple logger that is easy to understand and use. Only stdlib libraries are used.



## Example
Global clog

```go
package main

import (
	"log"

	"github.com/chrispassas/clog"
)

func main() {
    // Example here
    clog.Debugf("This is a log message")
}
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
