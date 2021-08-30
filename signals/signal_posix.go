package signals

import (
	"os"
	"syscall"
)

var shutdownSignals = []os.Signal{syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGHUP}
