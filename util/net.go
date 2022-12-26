package util

import (
	"encoding/binary"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/spf13/cast"
)

//////////// Network Functions ////////////

// GetHostName Gethostname gethostname()
func GetHostName() (string, error) {
	return os.Hostname()
}

// GetHostByName Gethostbyname gethostbyname()
// Get the IPv4 address corresponding to a given Internet host name
func GetHostByName(hostname string) (string, error) {
	ips, err := net.LookupIP(hostname)
	if ips != nil {
		for _, v := range ips {
			if v.To4() != nil {
				return v.String(), nil
			}
		}
		return "", nil
	}
	return "", err
}

// GetHostByNameL Gethostbynamel gethostbynamel()
// Get a list of IPv4 addresses corresponding to a given Internet host name
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

// GetHostByAddr Gethostbyaddr gethostbyaddr()
// Get the Internet host name corresponding to a given IP address
func GetHostByAddr(ipAddress string) (string, error) {
	names, err := net.LookupAddr(ipAddress)
	if names != nil {
		return strings.TrimRight(names[0], "."), nil
	}
	return "", err
}

// IP2long ip2long()
// IPv4
func IP2long(ipAddress string) uint32 {
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return 0
	}
	return binary.BigEndian.Uint32(ip.To4())
}

// Long2ip long2ip()
// IPv4
func Long2ip(properAddress uint32) string {
	ipByte := make([]byte, 4)
	binary.BigEndian.PutUint32(ipByte, properAddress)
	ip := net.IP(ipByte)
	return ip.String()
}

// ExtractIP returns a real ip
func ExtractIP(addr string) (string, error) {
	// if addr specified then its returned
	if len(addr) > 0 && (addr != "0.0.0.0" && addr != "[::]" && addr != "::") {
		return addr, nil
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("Failed to get interfaces! Err: %v", err)
	}

	var addrs []net.Addr
	var loAddrs []net.Addr
	for _, iface := range ifaces {
		ifaceAddrs, err := iface.Addrs()
		if err != nil {
			// ignore error, interface can dissapear from system
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

	// return private ip
	if ipAddr != nil {
		return net.IP(ipAddr).String(), nil
	}

	// return public or virtual ip
	if publicIP != nil {
		return net.IP(publicIP).String(), nil
	}

	return "", fmt.Errorf("No IP address found, and explicit IP not provided")
}

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

// URL wrap url.URL.
type URL struct {
	Scheme     string
	Opaque     string        // encoded opaque data
	User       *url.Userinfo // username and password information
	Host       string        // host or host:port
	Path       string        // path (relative paths may omit leading slash)
	RawPath    string        // encoded path hint (see EscapedPath method)
	ForceQuery bool          // append a query ('?') even if RawQuery is empty
	RawQuery   string        // encoded query values, without '?'
	Fragment   string        // fragment for references, without '#'
	HostName   string
	Port       string
	params     url.Values
}

// ParseURLRaw parses raw into URL.
func ParseURLRaw(raw string) (*URL, error) {
	u, e := url.Parse(raw)
	if e != nil {
		return nil, e
	}

	return &URL{
		Scheme:     u.Scheme,
		Opaque:     u.Opaque,
		User:       u.User,
		Host:       u.Host,
		Path:       u.Path,
		RawPath:    u.RawPath,
		ForceQuery: u.ForceQuery,
		RawQuery:   u.RawQuery,
		Fragment:   u.Fragment,
		HostName:   u.Hostname(),
		Port:       u.Port(),
		params:     u.Query(),
	}, nil
}

// Password gets password from URL.
func (u *URL) Password() (string, bool) {
	if u.User != nil {
		return u.User.Password()
	}
	return "", false
}

// Username gets username from URL.
func (u *URL) Username() string {
	return u.User.Username()
}

// QueryInt returns provided field's value in int type.
// if value is empty, expect returns
func (u *URL) QueryInt(field string, expect int) (ret int) {
	ret, err := cast.ToIntE(u.Query().Get(field))
	if err != nil {
		return expect
	}

	return ret
}

// QueryInt64 returns provided field's value in int64 type.
// if value is empty, expect returns
func (u *URL) QueryInt64(field string, expect int64) (ret int64) {
	ret, err := cast.ToInt64E(u.Query().Get(field))
	if err != nil {
		return expect
	}

	return ret
}

// QueryString returns provided field's value in string type.
// if value is empty, expect returns
func (u *URL) QueryString(field string, expect string) (ret string) {
	ret = expect
	if mi := u.Query().Get(field); mi != "" {
		ret = mi
	}

	return
}

// QueryDuration returns provided field's value in duration type.
// if value is empty, expect returns
func (u *URL) QueryDuration(field string, expect time.Duration) (ret time.Duration) {
	ret, err := cast.ToDurationE(u.Query().Get(field))
	if err != nil {
		return expect
	}

	return ret
}

// QueryBool returns provided field's value in bool
// if value is empty, expect returns
func (u *URL) QueryBool(field string, expect bool) (ret bool) {
	ret, err := cast.ToBoolE(u.Query().Get(field))
	if err != nil {
		return expect
	}
	return ret
}

// Query parses RawQuery and returns the corresponding values.
// It silently discards malformed value pairs.
// To check errors use ParseQuery.
func (u *URL) Query() url.Values {
	v, _ := url.ParseQuery(u.RawQuery)
	return v
}
