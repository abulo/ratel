package util

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Echo 将可变参数依次写入标准输出，参数之间不追加分隔符，也不自动换行。
func Echo(args ...any) {
	fmt.Print(args...)
}

// Uniqid 基于当前时间生成带前缀的唯一标识字符串，格式为 前缀 + 8 位 Unix 时间
// 戳十六进制 + 5 位纳秒余数十六进制。
func Uniqid(prefix string) string {
	now := time.Now()
	return fmt.Sprintf("%s%08x%05x", prefix, now.Unix(), now.UnixNano()%0x100000)
}

// Exit 以指定状态码终止当前进程。
func Exit(status int) {
	os.Exit(status)
}

// Die 以指定状态码终止当前进程，语义与 Exit 相同。
func Die(status int) {
	os.Exit(status)
}

// Getenv 读取环境变量 varname；变量不存在时返回空字符串。
func Getenv(varname string) string {
	return os.Getenv(varname)
}

// Putenv 按 "KEY=VALUE" 形式设置环境变量。
//
// setting 中必须恰好包含一个等号，否则 panic；调用系统设置失败时返回错误。
func Putenv(setting string) error {
	s := strings.Split(setting, "=")
	if len(s) != 2 {
		panic("setting: invalid")
	}
	return os.Setenv(s[0], s[1])
}

// MemoryGetUsage 返回当前进程已分配堆内存的字节数。
//
// realUsage 参数当前未被使用，返回值始终为运行时统计的 Alloc。
func MemoryGetUsage(realUsage bool) uint64 {
	stat := new(runtime.MemStats)
	runtime.ReadMemStats(stat)
	return stat.Alloc
}

// MemoryGetPeakUsage 返回进程启动以来累计分配内存的字节数。
//
// realUsage 参数当前未被使用，返回值始终为运行时统计的 TotalAlloc。
func MemoryGetPeakUsage(realUsage bool) uint64 {
	stat := new(runtime.MemStats)
	runtime.ReadMemStats(stat)
	return stat.TotalAlloc
}

// VersionCompare 按版本号规则比较 version1 与 version2，并按 operator 返回布尔结果。
//
// operator 支持 <、lt、<=、le、>、gt、>=、ge、==、=、eq、!=、<>、ne；非法运算符
// 会 panic。版本号中的预发布标识按以下顺序比较：普通字符串 < dev < alpha/a <
// beta/b < RC/rc < # < pl/p。
func VersionCompare(version1, version2, operator string) bool {
	var vcompare func(string, string) int
	var canonicalize func(string) string
	var special func(string, string) int

	vcompare = func(origV1, origV2 string) int {
		if origV1 == "" || origV2 == "" {
			if origV1 == "" && origV2 == "" {
				return 0
			}
			if origV1 == "" {
				return -1
			}
			return 1
		}

		ver1, ver2, compare := "", "", 0
		if origV1[0] == '#' {
			ver1 = origV1
		} else {
			ver1 = canonicalize(origV1)
		}
		if origV2[0] == '#' {
			ver2 = origV2
		} else {
			ver2 = canonicalize(origV2)
		}
		n1, n2 := 0, 0
		for {
			p1, p2 := "", ""
			n1 = strings.IndexByte(ver1, '.')
			if n1 == -1 {
				p1, ver1 = ver1[:], ""
			} else {
				p1, ver1 = ver1[:n1], ver1[n1+1:]
			}
			n2 = strings.IndexByte(ver2, '.')
			if n2 == -1 {
				p2, ver2 = ver2, ""
			} else {
				p2, ver2 = ver2[:n2], ver2[n2+1:]
			}

			if (p1[0] >= '0' && p1[0] <= '9') && (p2[0] >= '0' && p2[0] <= '9') {
				l1, _ := strconv.Atoi(p1)
				l2, _ := strconv.Atoi(p2)
				if l1 > l2 {
					compare = 1
				} else if l1 == l2 {
					compare = 0
				} else {
					compare = -1
				}
			} else if !(p1[0] >= '0' && p1[0] <= '9') && !(p2[0] >= '0' && p2[0] <= '9') {
				compare = special(p1, p2)
			} else {
				if p1[0] >= '0' && p1[0] <= '9' {
					compare = special("#N#", p2)
				} else {
					compare = special(p1, "#N#")
				}
			}

			if compare != 0 || n1 == -1 || n2 == -1 {
				break
			}
		}

		if compare == 0 {
			if ver1 != "" {
				if ver1[0] >= '0' && ver1[0] <= '9' {
					compare = 1
				} else {
					compare = vcompare(ver1, "#N#")
				}
			} else if ver2 != "" {
				if ver2[0] >= '0' && ver2[0] <= '9' {
					compare = -1
				} else {
					compare = vcompare("#N#", ver2)
				}
			}
		}

		return compare
	}

	canonicalize = func(version string) string {
		ver := []byte(version)
		l := len(ver)
		if l == 0 {
			return ""
		}
		var buf = make([]byte, l*2)
		j := 0
		for i, v := range ver {
			next := uint8(0)
			if i+1 < l {
				next = ver[i+1]
			}
			if v == '-' || v == '_' || v == '+' {
				if j > 0 && buf[j-1] != '.' {
					buf[j] = '.'
					j++
				}
			} else if (next > 0) &&
				(!(next >= '0' && next <= '9') && (v >= '0' && v <= '9')) ||
				(!(v >= '0' && v <= '9') && (next >= '0' && next <= '9')) {
				buf[j] = v
				j++
				if v != '.' && next != '.' {
					buf[j] = '.'
					j++
				}
				continue
			} else if !((v >= '0' && v <= '9') ||
				(v >= 'a' && v <= 'z') || (v >= 'A' && v <= 'Z')) {
				if j > 0 && buf[j-1] != '.' {
					buf[j] = '.'
					j++
				}
			} else {
				buf[j] = v
				j++
			}
		}

		return string(buf[:j])
	}

	special = func(form1, form2 string) int {
		found1, found2, len1, len2 := -1, -1, len(form1), len(form2)
		// 预发布标识的先后顺序：未识别字符串 < dev < alpha/a < beta/b < RC/rc < # < pl/p。
		forms := map[string]int{
			"dev":   0,
			"alpha": 1,
			"a":     1,
			"beta":  2,
			"b":     2,
			"RC":    3,
			"rc":    3,
			"#":     4,
			"pl":    5,
			"p":     5,
		}

		for name, order := range forms {
			if len1 < len(name) {
				continue
			}
			if strings.Compare(form1[:len(name)], name) == 0 {
				found1 = order
				break
			}
		}
		for name, order := range forms {
			if len2 < len(name) {
				continue
			}
			if strings.Compare(form2[:len(name)], name) == 0 {
				found2 = order
				break
			}
		}

		if found1 == found2 {
			return 0
		} else if found1 > found2 {
			return 1
		} else {
			return -1
		}
	}

	compare := vcompare(version1, version2)

	switch operator {
	case "<", "lt":
		return compare == -1
	case "<=", "le":
		return compare != 1
	case ">", "gt":
		return compare == 1
	case ">=", "ge":
		return compare != -1
	case "==", "=", "eq":
		return compare == 0
	case "!=", "<>", "ne":
		return compare != 0
	default:
		panic("operator: invalid")
	}
}

// ZipOpen 打开 filename 指向的 zip 文件并返回读取器；调用方负责在使用完毕后关闭读取器。
func ZipOpen(filename string) (*zip.ReadCloser, error) {
	return zip.OpenReader(filename)
}

// Pack 按 order 指定的字节序将 data 序列化为二进制缓冲区并以字符串返回；写入失败时
// 返回空字符串与错误。
func Pack(order binary.ByteOrder, data any) (string, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, order, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// Unpack 按 order 指定的字节序从二进制字符串 data 反序列化；读取失败时返回 nil 与错误。
func Unpack(order binary.ByteOrder, data string) (any, error) {
	var result []byte
	r := bytes.NewReader([]byte(data))
	err := binary.Read(r, order, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Ternary 在 condition 为真时返回 trueVal，否则返回 falseVal。
//
// 两个返回值为 any 类型，调用方需要自行做类型断言。
func Ternary(condition bool, trueVal, falseVal any) any {
	if condition {
		return trueVal
	}
	return falseVal
}
