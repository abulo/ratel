package env

var (
	appMode     string
	appRegion   string
	appZone     string
	appHost     string
	appInstance string
	appName     string
	appID       string
	startTime   string
)

//StartTime get start time
func StartTime() string {
	return startTime
}

//SetStartTime set StartTime
func SetStartTime(s string) {
	startTime = s
}

//AppID get appID
func AppID() string {
	return appID
}

//SetAppID set appID
func SetAppID(s string) {
	appID = s
}

// Name gets application name.
func Name() string {
	return appName
}

//SetName set app anme
func SetName(s string) {
	appName = s
}

func AppMode() string {
	return appMode
}

func SetAppMode(appMode string) {
	appMode = appMode
}

func AppRegion() string {
	return appRegion
}

func SetAppRegion(appRegion string) {
	appRegion = appRegion
}

func AppZone() string {
	return appZone
}

func SetAppZone(appZone string) {
	appZone = appZone
}

func AppHost() string {
	return appHost
}

func SetAppHost(appHost string) {
	appHost = appHost
}

func AppInstance() string {
	return appInstance
}
