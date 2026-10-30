package util

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

// GetHostName 返回本机主机名；系统获取失败时返回错误。
func GetHostName() (string, error) {
	return os.Hostname()
}

// GetHostByName 返回 hostname 对应的第一个 IPv4 地址。
//
// DNS 查询成功但结果中不包含 IPv4 地址时返回错误；查询本身失败时返回底层错误。
func GetHostByName(hostname string) (string, error) {
	ips, err := net.LookupIP(hostname)
	if ips != nil {
		for _, v := range ips {
			if v.To4() != nil {
				return v.String(), nil
			}
		}
		return "", fmt.Errorf("null: cannot find host %s", hostname)
	}
	return "", err
}

// GetHostByNameL 返回 hostname 对应的全部 IPv4 地址列表；查询失败时返回错误。
func GetHostByNameL(hostname string) ([]string, error) {
	ips, err := net.LookupIP(hostname)
	if ips != nil {
		var ipstrs []string
		for _, v := range ips {
			if v.To4() != nil {
				ipstrs = append(ipstrs, v.String())
			}
		}
		return ipstrs, nil
	}
	return nil, err
}

// GetHostByAddr 返回 IP 地址 ipAddress 反向解析得到的主机名，末尾的点号被去除；
// 解析失败时返回错误。
func GetHostByAddr(ipAddress string) (string, error) {
	names, err := net.LookupAddr(ipAddress)
	if names != nil {
		return strings.TrimRight(names[0], "."), nil
	}
	return "", err
}

// IP2long 将 IPv4 地址转换为大端 32 位无符号整数；地址非法时返回 0。
func IP2long(ipAddress string) uint32 {
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return 0
	}
	return binary.BigEndian.Uint32(ip.To4())
}

// Long2ip 将大端 32 位无符号整数还原为 IPv4 地址字符串。
func Long2ip(properAddress uint32) string {
	ipByte := make([]byte, 4)
	binary.BigEndian.PutUint32(ipByte, properAddress)
	ip := net.IP(ipByte)
	return ip.String()
}

// ExtractIP 解析可用的本机 IP 地址。
//
// addr 非空且不是 0.0.0.0、:: 等通配地址时直接返回；否则遍历本机网卡，优先返回
// 内网地址，内网不存在时再返回公网或虚拟地址；回环地址排在普通网卡之后尝试。
// 找不到任何地址时返回错误。
func ExtractIP(addr string) (string, error) {
	if len(addr) > 0 && (addr != "0.0.0.0" && addr != "[::]" && addr != "::") {
		return addr, nil
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("failed to get interfaces! Err: %v", err)
	}

	var addrs []net.Addr
	var loAddrs []net.Addr
	for _, iface := range ifaces {
		ifaceAddrs, err := iface.Addrs()
		if err != nil {
			// 网卡可能在运行过程中消失，跳过本次错误即可。
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
			loAddrs = append(loAddrs, ifaceAddrs...)
			continue
		}
		addrs = append(addrs, ifaceAddrs...)
	}
	addrs = append(addrs, loAddrs...)

	var ipAddr []byte
	var publicIP []byte

	for _, rawAddr := range addrs {
		var ip net.IP
		switch addr := rawAddr.(type) {
		case *net.IPAddr:
			ip = addr.IP
		case *net.IPNet:
			ip = addr.IP
		default:
			continue
		}

		if !isPrivateIP(ip.String()) {
			publicIP = ip
			continue
		}

		ipAddr = ip
		break
	}

	if ipAddr != nil {
		return net.IP(ipAddr).String(), nil
	}

	if publicIP != nil {
		return net.IP(publicIP).String(), nil
	}

	return "", errors.New("no IP address found, and explicit IP not provided")
}

// isPrivateIP 判断 ipAddr 是否属于内网、保留或运营商 NAT 保留网段。
func isPrivateIP(ipAddr string) bool {
	var privateBlocks []*net.IPNet

	for _, b := range []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "100.64.0.0/10", "fd00::/8"} {
		if _, block, err := net.ParseCIDR(b); err == nil {
			privateBlocks = append(privateBlocks, block)
		}
	}

	ip := net.ParseIP(ipAddr)
	for _, priv := range privateBlocks {
		if priv.Contains(ip) {
			return true
		}
	}
	return false
}
