## Quick Start
Generate an ID
```go

import (
  "fmt"
  "github.com/abulo/ratel/v3/snowflake"
)

func ExampleGenInt64ID() {
  id := snowflake.CommonConfig.GenInt64ID()
  fmt.Printf("id generated: %v", id)
}
```