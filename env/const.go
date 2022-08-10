package env

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/abulo/ratel/v3/constant"
	"github.com/abulo/ratel/v3/util"
)

var (
	appMode     string
	appRegion   string
	appZone     string
	appHost     string
	appInstance string
	appPodIP    string
	appPodName  string
)

func InitEnv() {
	appID = os.Getenv(constant.EnvAppID)
	appMode = os.Getenv(constant.EnvAppMode)
	appRegion = os.Getenv(constant.EnvAppRegion)
	appZone = os.Getenv(constant.EnvAppZone)
	appHost = os.Getenv(constant.EnvAppHost)
	appInstance = os.Getenv(constant.EnvAppInstance)
	if appInstance == "" {
		appInstance = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s:%s", HostName(), AppID()))))
	}
	appPodIP = os.Getenv(constant.EnvPOD_IP)
	appPodName = os.Getenv(constant.EnvPOD_NAME)
}

func AppMode() string {
	return appMode
}

func SetAppMode(mode string) {
	appMode = mode
}

func AppRegion() string {
	return appRegion
}

func SetAppRegion(region string) {
	appRegion = region
}

func AppZone() string {
	return appZone
}

func SetAppZone(zone string) {
	appZone = zone
}

func AppHost() string {
	return appHost
}

func SetAppHost(host string) {
	appHost = host
}

func AppInstance() string {
	return appInstance
}

func SetAppInstance(instance string) {
	appInstance = instance
}

const ratelVersion = "ratel3.1.5"

var (
	startTime string
	goVersion string
)

// build info
/*

 */
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

// Name gets application name.
func Name() string {
	return appName
}

// SetName set app name
func SetName(s string) {
	appName = s
}

// AppID get appID
func AppID() string {
	if appID == "" {
		return "1234567890" //default appid when APP_ID Env var not set
	}
	return appID
}

// SetAppID set appID
func SetAppID(s string) {
	appID = s
}

// AppVersion get buildAppVersion
func BuildVersion() string {
	return buildVersion
}

// SetAppVersion set buildVersion
func SetBuildVersion(s string) {
	buildVersion = s
}

// RatelVersion get ratelVersion
func RatelVersion() string {
	return ratelVersion
}

// BuildTime get buildTime
func BuildTime() string {
	return buildTime
}

// SetBuildTime set buildTime
func SetBuildTime(param string) {
	buildTime = param
}

// HostName get host name
func HostName() string {
	return hostName
}

// StartTime get start time
func StartTime() string {
	return startTime
}

// GoVersion get go version
func GoVersion() string {
	return goVersion
}

// PrintVersion print format version info
func PrintVersion() {
	fmt.Printf("%-8s]> %-30s => %s\n", "ratel", "GoVersion", GoVersion())
	fmt.Printf("%-8s]> %-30s => %s\n", "ratel", "RatelVersion", ratelVersion)
	fmt.Printf("%-8s]> %-30s => %s\n", "ratel", "Name", Name())
	fmt.Printf("%-8s]> %-30s => %s\n", "ratel", "AppID", AppID())
	fmt.Printf("%-8s]> %-30s => %s\n", "ratel", "Region", AppRegion())
	fmt.Printf("%-8s]> %-30s => %s\n", "ratel", "Zone", AppZone())
	fmt.Printf("%-8s]> %-30s => %s\n", "ratel", "BuildVersion", BuildVersion())
	fmt.Printf("%-8s]> %-30s => %s\n", "ratel", "BuildTime", BuildTime())
	fmt.Printf("%-8s]> %-30s => %s\n", "ratel", "StartTime", StartTime())
}
