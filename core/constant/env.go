package constant

// EnvKeySentinelLogDir ...
const (
	// EnvKeySentinelLogDir ...
	EnvKeySentinelLogDir = "SENTINEL_LOG_DIR"
	// EnvKeySentinelAppName ...
	EnvKeySentinelAppName = "SENTINEL_APP_NAME"
)

// EnvAppName ...
const (
	// EnvAppName ...
	EnvAppName     = "APP_NAME"
	EnvAppID       = "APP_ID"
	EnvDeployment  = "APP_DEPLOYMENT"
	EnvAppMode     = "APP_MODE"
	EnvAppRegion   = "APP_REGION"
	EnvAppZone     = "APP_ZONE"
	EnvAppHost     = "APP_HOST"
	EnvAppInstance = "APP_INSTANCE" // application unique instance id.
	EnvPodIP       = "POD_IP"       //k8s环境
	EnvPodNAME     = "POD_NAME"     //k8s环境
)

// DefaultDeployment ...
const (
	// DefaultDeployment ...
	DefaultDeployment = ""
	// DefaultRegion ...
	DefaultRegion = ""
	// DefaultZone ...
	DefaultZone = ""
)

// KeyBalanceGroup ...
const (
	// KeyBalanceGroup ...
	KeyBalanceGroup = "__group"

	// DefaultBalanceGroup ...
	DefaultBalanceGroup = "default"
)
