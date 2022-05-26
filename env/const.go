package env

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

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

const ratelVersion = "3.0.3"

var (
	startTime string
	goVersion string
)

// build info
/*

 */
var (
	appName         string
	appID           string
	hostName        string
	buildAppVersion string
	buildUser       string
	buildHost       string
	buildStatus     string
	buildTime       string
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

//SetName set app name
func SetName(s string) {
	appName = s
}

//AppID get appID
func AppID() string {
	if appID == "" {
		return "1234567890" //default appid when APP_ID Env var not set
	}
	return appID
}

//SetAppID set appID
func SetAppID(s string) {
	appID = s
}

//AppVersion get buildAppVersion
func AppVersion() string {
	return buildAppVersion
}

// SetAppVersion set appVersion
func SetAppVersion(s string) {
	buildAppVersion = s
}

//RatelVersion get ratelVersion
func RatelVersion() string {
	return ratelVersion
}

//BuildTime get buildTime
func BuildTime() string {
	return buildTime
}

//BuildUser get buildUser
func BuildUser() string {
	return buildUser
}

//BuildHost get buildHost
func BuildHost() string {
	return buildHost
}

//SetBuildTime set buildTime
func SetBuildTime(param string) {
	buildTime = strings.Replace(param, "--", " ", 1)
}

// HostName get host name
func HostName() string {
	return hostName
}

//StartTime get start time
func StartTime() string {
	return startTime
}

//GoVersion get go version
func GoVersion() string {
	return goVersion
}

// PrintVersion print format version info
func PrintVersion() {
	fmt.Printf("%-8s]> %-30s => %s\n", "ratel", "name", appName)
	fmt.Printf("%-8s]> %-30s => %s\n", "ratel", "appID", appID)
	fmt.Printf("%-8s]> %-30s => %s\n", "ratel", "region", AppRegion())
	fmt.Printf("%-8s]> %-30s => %s\n", "ratel", "zone", AppZone())
	fmt.Printf("%-8s]> %-30s => %s\n", "ratel", "appVersion", buildAppVersion)
	fmt.Printf("%-8s]> %-30s => %s\n", "ratel", "ratelVersion", ratelVersion)
	fmt.Printf("%-8s]> %-30s => %s\n", "ratel", "buildUser", buildUser)
	fmt.Printf("%-8s]> %-30s => %s\n", "ratel", "buildHost", buildHost)
	fmt.Printf("%-8s]> %-30s => %s\n", "ratel", "buildTime", BuildTime())
	fmt.Printf("%-8s]> %-30s => %s\n", "ratel", "buildStatus", buildStatus)
}