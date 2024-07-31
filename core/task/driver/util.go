package driver

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// GlobalKeyPrefix is a global redis key prefix
const GlobalKeyPrefix = "Distributed-Cron:"

var KeyPreFix string

func GetKeyPre(serviceName string) string {
	return GlobalKeyPrefix + serviceName + ":"
}

func SetNodeIdKeyPrefix(prefix string) {
	if prefix != "" {
		KeyPreFix = strings.Replace(prefix, "%s", "", -1)
	}
}

func GetNodeIdKeyPrefix() string {
	return KeyPreFix
}

func GetNodeId(serviceName string) string {
	return GetKeyPre(serviceName) + uuid.New().String()
}

func TimePre(t time.Time, preDuration time.Duration) int64 {
	return t.Add(-preDuration).Unix()
}
