package constant

const (
	// EnvKeySentinelLogDir ...
	EnvKeySentinelLogDir = "SENTINEL_LOG_DIR"
	// EnvKeySentinelAppName ...
	EnvKeySentinelAppName = "SENTINEL_APP_NAME"
)

const (
	// EnvAppName ...
	EnvAppName = "APP_NAME"
	EnvAppID   = "APP_ID"
	// EnvDeployment ...
	EnvDeployment = "APP_DEPLOYMENT"

	EnvAppLogDir   = "APP_LOG_DIR"
	EnvAppMode     = "APP_MODE"
	EnvAppRegion   = "APP_REGION"
	EnvAppZone     = "APP_ZONE"
	EnvAppHost     = "APP_HOST"
	EnvAppInstance = "APP_INSTANCE" // application unique instance id.
	//k8s环境
	EnvPOD_IP   = "POD_IP"
	EnvPOD_NAME = "POD_NAME"
)

const (
	// DefaultDeployment ...
	DefaultDeployment = ""
	// DefaultRegion ...
	DefaultRegion = ""
	// DefaultZone ...
	DefaultZone = ""
)

const (
	// KeyBalanceGroup ...
	KeyBalanceGroup = "__group"

	// DefaultBalanceGroup ...
	DefaultBalanceGroup = "default"
)
