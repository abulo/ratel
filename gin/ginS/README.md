# Gin Default Server

This is API experiment for Gin.

```go
package main

import (
	"github.com/abulo/ratel/gin"
	"github.com/abulo/ratel/gin/ginS"
)

func main() {
	ginS.GET("/", func(c *gin.Context) { c.String(200, "Hello World") })
	ginS.Run()
}
```
