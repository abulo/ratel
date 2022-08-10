package entry

import (
	"time"
)

// Entry ...
type Entry struct {
	Host      string                 `json:"host"`
	Timestamp time.Time              `json:"timestamp"`
	File      string                 `json:"file"`
	Func      string                 `json:"func"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data"`
	Level     string                 `json:"level"`
}
