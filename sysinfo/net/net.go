package net

import "github.com/shirou/gopsutil/net"

// Useage Total总量，Idle空闲，Used使用率，Collercter总量，使用量
type Useage struct {
	BytesRecv uint64 `json:"recv"`
	BytesSent uint64 `json:"send"`
	Name      string `json:"name"`
}

// GetInfo 获取磁盘使用信息
func GetInfo() (useages []*Useage, err error) {
	nv, err := net.IOCounters(false)
	if err != nil {
		return
	}
	useages = make([]*Useage, 0, len(nv))
	for _, status := range nv {
		useage := &Useage{BytesRecv: status.BytesRecv, BytesSent: status.BytesSent, Name: status.Name}
		useages = append(useages, useage)
	}
	return
}
