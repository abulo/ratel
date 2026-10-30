package hostname

import "os"

// Hostname 返回当前主机名；获取主机名的系统调用失败时返回 "unknown"。
func Hostname() string {
	hostName, err := os.Hostname()
	if err != nil {
		hostName = "unknown"
	}
	return hostName
}
