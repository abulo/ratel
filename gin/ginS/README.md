# Gin Default Server

This is API experiment for Gin.

```go
package main

import (
	"github.com/abulo/ratel/v2/gin"
	"github.com/abulo/ratel/v2/gin/ginS"
)

func main() {
	ginS.GET("/", func(c *gin.Context) { c.String(200, "Hello World") })
	ginS.Run()
}
```
