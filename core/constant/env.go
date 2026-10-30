package constant

// Sentinel 组件读取的环境变量名常量。
const (
	// EnvKeySentinelLogDir 为 Sentinel 规则与日志目录的环境变量名。
	EnvKeySentinelLogDir = "SENTINEL_LOG_DIR"
	// EnvKeySentinelAppName 为 Sentinel 应用名称的环境变量名。
	EnvKeySentinelAppName = "SENTINEL_APP_NAME"
)

// 应用元信息环境变量名常量，由部署环境注入，供服务发现与日志记录使用。
const (
	// EnvAppName 为应用名称的环境变量名。
	EnvAppName = "APP_NAME"
	// EnvAppID 为应用编号的环境变量名。
	EnvAppID = "APP_ID"
	// EnvDeployment 为部署环境标识的环境变量名。
	EnvDeployment = "APP_DEPLOYMENT"
	// EnvAppMode 为运行模式的环境变量名。
	EnvAppMode = "APP_MODE"
	// EnvAppRegion 为部署地域的环境变量名。
	EnvAppRegion = "APP_REGION"
	// EnvAppZone 为部署可用区的环境变量名。
	EnvAppZone = "APP_ZONE"
	// EnvAppHost 为运行主机的环境变量名。
	EnvAppHost = "APP_HOST"
	// EnvAppInstance 为应用唯一实例标识的环境变量名。
	EnvAppInstance = "APP_INSTANCE"
	// EnvPodIP 为 Kubernetes 环境中 Pod IP 的环境变量名。
	EnvPodIP = "POD_IP"
	// EnvPodNAME 为 Kubernetes 环境中 Pod 名称的环境变量名。
	EnvPodNAME = "POD_NAME"
)

// 应用部署元信息缺省值，环境变量未提供时使用。
const (
	// DefaultDeployment 为部署环境标识的默认值，空字符串表示未指定。
	DefaultDeployment = ""
	// DefaultRegion 为部署地域的默认值，空字符串表示未指定。
	DefaultRegion = ""
	// DefaultZone 为部署可用区的默认值，空字符串表示未指定。
	DefaultZone = ""
)

// 负载均衡分组相关常量。
const (
	// KeyBalanceGroup 为分组元信息使用的键名。
	KeyBalanceGroup = "__group"

	// DefaultBalanceGroup 为未显式指定分组时的默认分组名。
	DefaultBalanceGroup = "default"
)
