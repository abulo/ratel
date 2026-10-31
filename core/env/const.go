package env

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/abulo/ratel/v3/core/constant"
	"github.com/abulo/ratel/v3/util"
)

const (
	ratelVersion = "v3.8.56"
)

var (
	appMode     string
	appRegion   string
	appZone     string
	appHost     string
	appInstance string
	// appPodIP    string
	// appPodName  string
)

// InitEnv 从环境变量加载应用模式、地域、可用区、主机与实例标识；
// 实例标识为空时，按主机名与 AppID 拼接后的 MD5 摘要生成唯一实例标识。
func InitEnv() {
	appID = os.Getenv(constant.EnvAppID)
	appMode = os.Getenv(constant.EnvAppMode)
	appRegion = os.Getenv(constant.EnvAppRegion)
	appZone = os.Getenv(constant.EnvAppZone)
	appHost = os.Getenv(constant.EnvAppHost)
	appInstance = os.Getenv(constant.EnvAppInstance)
	if appInstance == "" {
		appInstance = fmt.Sprintf("%x", md5.Sum(fmt.Appendf(nil, "%s:%s", HostName(), AppID())))
	}
	// appPodIP = os.Getenv(constant.EnvPOD_IP)
	// appPodName = os.Getenv(constant.EnvPOD_NAME)
}

// AppMode 返回当前运行模式，取值来自环境变量 APP_MODE。
func AppMode() string {
	return appMode
}

// SetAppMode 覆盖当前运行模式。
func SetAppMode(mode string) {
	appMode = mode
}

// AppRegion 返回应用所在地域，取值来自环境变量 APP_REGION。
func AppRegion() string {
	return appRegion
}

// SetAppRegion 覆盖应用所在地域。
func SetAppRegion(region string) {
	appRegion = region
}

// AppZone 返回应用所在可用区，取值来自环境变量 APP_ZONE。
func AppZone() string {
	return appZone
}

// SetAppZone 覆盖应用所在可用区。
func SetAppZone(zone string) {
	appZone = zone
}

// AppHost 返回应用运行主机标识，取值来自环境变量 APP_HOST。
func AppHost() string {
	return appHost
}

// SetAppHost 覆盖应用运行主机标识。
func SetAppHost(host string) {
	appHost = host
}

// AppInstance 返回应用实例唯一标识。
func AppInstance() string {
	return appInstance
}

// SetAppInstance 覆盖应用实例唯一标识。
func SetAppInstance(instance string) {
	appInstance = instance
}

var (
	startTime string
	goVersion string
)

var (
	appName      string
	appID        string
	hostName     string
	buildVersion string
	buildTime    string
)

func init() {
	if appName == "" {
		appName = os.Getenv(constant.EnvAppName)
		if appName == "" {
			appName = filepath.Base(os.Args[0])
		}
	}

	name, err := os.Hostname()
	if err != nil {
		name = "unknown"
	}
	hostName = name
	startTime = util.Date("Y-m-d H:i:s", util.Now())
	SetBuildTime(buildTime)
	goVersion = runtime.Version()
	InitEnv()
}

// Name 返回应用名称；优先取环境变量 APP_NAME，未设置时取可执行文件名。
func Name() string {
	return appName
}

// SetName 覆盖应用名称。
func SetName(s string) {
	appName = s
}

// AppID 返回应用 ID；环境变量 APP_ID 未设置时返回默认值 "1234567890"。
func AppID() string {
	if appID == "" {
		return "1234567890"
	}
	return appID
}

// SetAppID 覆盖应用 ID。
func SetAppID(s string) {
	appID = s
}

// BuildVersion 返回构建版本号，由外部在构建或启动阶段通过 SetBuildVersion 写入。
func BuildVersion() string {
	return buildVersion
}

// SetBuildVersion 覆盖构建版本号。
func SetBuildVersion(s string) {
	buildVersion = s
}

// RatelVersion 返回框架当前版本号。
func RatelVersion() string {
	return ratelVersion
}

// BuildTime 返回构建时间，由外部在构建或启动阶段通过 SetBuildTime 写入。
func BuildTime() string {
	return buildTime
}

// SetBuildTime 覆盖构建时间。
func SetBuildTime(param string) {
	buildTime = param
}

// HostName 返回运行所在主机名；获取失败时返回 "unknown"。
func HostName() string {
	return hostName
}

// StartTime 返回应用本次启动时间，格式为 Y-m-d H:i:s。
func StartTime() string {
	return startTime
}

// GoVersion 返回编译当前程序使用的 Go 版本。
func GoVersion() string {
	return goVersion
}

// PrintVersion 将应用名称、ID、版本、地域、可用区、主机与启动时间等环境信息格式化输出到标准输出。
func PrintVersion() {
	fmt.Printf("%-8s]> %-30s => %s\n", "App", "GoVersion", GoVersion())
	fmt.Printf("%-8s]> %-30s => %s\n", "App", "RatelVersion", RatelVersion())
	fmt.Printf("%-8s]> %-30s => %s\n", "App", "Name", Name())
	fmt.Printf("%-8s]> %-30s => %s\n", "App", "AppID", AppID())
	fmt.Printf("%-8s]> %-30s => %s\n", "App", "Region", AppRegion())
	fmt.Printf("%-8s]> %-30s => %s\n", "App", "Zone", AppZone())
	fmt.Printf("%-8s]> %-30s => %s\n", "App", "BuildVersion", BuildVersion())
	fmt.Printf("%-8s]> %-30s => %s\n", "App", "BuildTime", BuildTime())
	fmt.Printf("%-8s]> %-30s => %s\n", "App", "StartTime", StartTime())
}
