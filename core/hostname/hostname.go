package hostname

import "os"

func Hostname() string {
	// Get the hostname
	hostName, err := os.Hostname()
	if err != nil {
		hostName = "unknown"
	}
	return hostName
}
