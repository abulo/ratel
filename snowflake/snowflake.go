package snowflake

import (
	"net"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
)

const (
	// StartTimeEnvName 为配置 ID 生成起始纪元的环境变量名。
	StartTimeEnvName = "SNOWFLAKEID_START_TIME"
	// MachineIDEnvName 为配置机器编号的环境变量名。
	MachineIDEnvName           = "SNOWFLAKEID_MACHINE_ID"
	defaultStartTimeStr        = "2018-06-01T00:00:00.000Z"
	defaultStartTimeNano int64 = 1527811200000000000
)

// Config 保存雪花 ID 的位布局、各段掩码、机器编号与上一次生成的 ID。
//
// 零值不可直接使用，须通过 ConstructConfig 或 ConstructConfigWithMachineID 构造；
// GenInt64ID 通过原子操作更新 lastID，可在并发环境下安全使用。
type Config struct {
	timeBits      uint
	seqBits       uint
	machineBits   uint
	timeMask      int64
	seqMask       int64
	machineID     int64
	machineIDMask int64
	lastID        int64
}

// ConstructConfig 按给定的时间位、序列号位、机器编号位宽度构造 Config，
// 机器编号由环境变量或本机 IP 自动探测得到。
func ConstructConfig(timeBits, seqBits, machineBits uint) *Config {
	return ConstructConfigWithMachineID(timeBits, seqBits, machineBits, getMachineID())
}

// ConstructConfigWithMachineID 按给定位宽度与显式机器编号构造 Config；
// machineID 会按机器位宽度截断，避免越界占用其他位段。
func ConstructConfigWithMachineID(timeBits, seqBits, machineBits uint, machineID int64) *Config {
	machineIDMask := ^(int64(-1) << machineBits)
	return &Config{
		timeBits:      timeBits,
		seqBits:       seqBits,
		machineBits:   machineBits,
		timeMask:      ^(int64(-1) << timeBits),
		seqMask:       ^(int64(-1) << seqBits),
		machineIDMask: machineIDMask,
		machineID:     machineID & machineIDMask,
		lastID:        0,
	}
}

// BenchmarkConfig 面向性能基准场景的位布局：40 位时间戳、15 位序列号、8 位机器编号。
var BenchmarkConfig = ConstructConfig(40, 15, 8)

// CommonConfig 面向一般应用场景的位布局：40 位时间戳、7 位序列号、16 位机器编号。
var CommonConfig = ConstructConfig(40, 7, 16)

var startEpochNano = getStartEpochFromEnv()

// getCurrentTimestamp 返回距起始纪元的当前时间戳段，并按时间位掩码截断。
func (c *Config) getCurrentTimestamp() int64 {
	// 纳秒右移 20 位换算为毫秒量级，再按 timeMask 截断到时间位宽度。
	return (time.Now().UnixNano() - startEpochNano) >> 20 & c.timeMask
}

// GenInt64ID 按当前 Config 的位布局生成一个唯一 int64 ID，可并发调用。
//
// 同一时间片内序列号耗尽时会自旋等待到下一个时间片；并发抢占更新 lastID 失败时
// 通过 CAS 重试，保证同一 Config 实例不会生成重复 ID。
func (c *Config) GenInt64ID() int64 {
	// PERF-7-36 [P1] CAS 冲突先本地短自旋 3 次再进入睡眠退避；单次循环只读取一次系统时钟并向下传递给序列号耗尽分支复用（冲突率与睡眠影响需并发 benchmark 实测验证，预期收益需实测验证）
	const maxLocalSpin = 3
	conflict := 0
	for {
		localLastID := atomic.LoadInt64(&c.lastID)
		seq := c.GetSeqFromID(localLastID)
		lastIDTime := c.GetTimeFromID(localLastID)
		rawNano := time.Now().UnixNano()
		now := (rawNano - startEpochNano) >> 20 & c.timeMask
		if now > lastIDTime {
			seq = 0
		} else if seq >= c.seqMask {
			time.Sleep(time.Duration(0xFFFFF - (rawNano & 0xFFFFF)))
			continue
		} else {
			seq++
		}

		newID := now<<(c.machineBits+c.seqBits) + seq<<c.machineBits + c.machineID
		if atomic.CompareAndSwapInt64(&c.lastID, localLastID, newID) {
			return newID
		}
		conflict++
		if conflict <= maxLocalSpin {
			// 本地短自旋直接重试，避免前几次冲突就进入协程睡眠/调度
			continue
		}
		time.Sleep(time.Duration(20))
	}
}

// GetSeqFromID 从已有 ID 中提取序列号段。
func (c *Config) GetSeqFromID(id int64) int64 {
	return (id >> c.machineBits) & c.seqMask
}

// GetTimeFromID 从已有 ID 中提取时间戳段。
func (c *Config) GetTimeFromID(id int64) int64 {
	return id >> (c.machineBits + c.seqBits)
}

// getMachineID 依次从环境变量 MachineIDEnvName、本机内网 IPv4 后两段获取机器编号；
// 两者均不可用时返回 0。
func getMachineID() int64 {
	if machineIDStr, ok := os.LookupEnv(MachineIDEnvName); ok {
		if machineID, err := strconv.ParseInt(machineIDStr, 10, 64); err == nil {
			return machineID
		}
	}
	// 环境变量未配置时，取内网 IPv4 的后两段作为机器编号。
	if ip, err := getIP(); err == nil {
		return (int64(ip[2]) << 8) + int64(ip[3])
	}
	return 0
}

// getStartEpochFromEnv 返回 ID 生成使用的起始纪元纳秒时间戳。
//
// TODO(owner): 现状为解析成功时返回内置默认纪元、解析失败时返回解析时间，
// 环境变量覆盖语义与预期相反，待后续修正。
func getStartEpochFromEnv() int64 {
	startTimeStr := getEnv(StartTimeEnvName, defaultStartTimeStr)
	var startEpochTime, err = time.Parse(time.RFC3339, startTimeStr)

	if err == nil {
		return defaultStartTimeNano
	}

	return startEpochTime.UnixNano()
}

// getIP 返回本机一张非回环的私有 IPv4 地址；找不到时返回错误。
func getIP() (net.IP, error) {
	if addrs, err := net.InterfaceAddrs(); err == nil {
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok {
				if !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
					ip := ipNet.IP.To4()

					if ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168 {
						return ip, nil
					}
				}
			}
		}
	}
	return nil, errors.New("Failed to get ip address")
}

// getEnv 读取环境变量，未设置时返回 fallback。
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
