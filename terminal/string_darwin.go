//go:build darwin
// +build darwin

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
	return sprint(YellowColor, msg, arg...)
}

// Red ...
func Red(msg string, arg ...interface{}) string {
	return sprint(RedColor, msg, arg...)
}

// Blue ...
func Blue(msg string, arg ...interface{}) string {
	return sprint(BlueColor, msg, arg...)
}

// Green ...
func Green(msg string, arg ...interface{}) string {
	return sprint(GreenColor, msg, arg...)
}

// Greenf ...
func Greenf(msg string, arg ...interface{}) string {
	return sprint(GreenColor, msg, arg...)
}

// sprint
func sprint(colorValue int, msg string, arg ...interface{}) string {
	if arg != nil {
		return fmt.Sprintf("\x1b[%dm%s\x1b[0m %+v", colorValue, msg, arrToTransform(arg))
	} else {
		return fmt.Sprintf("\x1b[%dm%s\x1b[0m", colorValue, msg)
	}
}
