//go:build windows
// +build windows

package terminal

import (
	"fmt"
	"math/rand"
	"strconv"
)

var _ = RandomColor()

// RandomColor generates a random color.
func RandomColor() string {
	return fmt.Sprintf("#%s", strconv.FormatInt(int64(rand.Intn(16777216)), 16))
}

// Yellow ...
func Yellow(msg string, arg ...interface{}) string {
	return sprint(msg, arg...)
}

// Red ...
func Red(msg string, arg ...interface{}) string {
	return sprint(msg, arg...)
}

// Blue ...
func Blue(msg string, arg ...interface{}) string {
	return sprint(msg, arg...)
}

// Green ...
func Green(msg string, arg ...interface{}) string {
	return sprint(msg, arg...)
}

// Greenf ...
func Greenf(msg string, arg ...interface{}) string {
	return sprint(msg, arg...)
}

// sprint ...
func sprint(msg string, arg ...interface{}) string {
	if arg != nil {
		return fmt.Sprintf("%s %+v\n", msg, arrToTransform(arg))
	}
	return fmt.Sprintf("%s", msg)
}
