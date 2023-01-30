## log 库
实现了简单了日志记录功能

**快速使用**

```go
package main

import "logging"

var log *logging.Logger

func main() {
	log = logging.NewLogger(logging.DEBUG)
	log.Debug("这是一条Debug信息")
	
	log.SetFile("path", 0644, 128*8)
	log.Warning("这是一条Warning信息")
}
```