package redis

import (
	"os"

	"github.com/spf13/viper"
)

const (
	// DefaultDir是默认配置文件搜索dir
	DefaultDir = "/etc/paas/"
	// DefaultFileName是默认的配置文件名称
	DefaultFileName = "redis"
	// efaultEnvPrefixKey是环境变量的缺省前缀
	DefaultEnvPrefixKey = ""
	// EnvPrefixKey是环境变量的前缀
	EnvPrefixKey = "ENV_PREFIX"
	// ConfigDirKey是关键文件搜索路径环境变量
	ConfigDirKey = "CONFIG_DIR"
	// ConfigNameKey从环境变量获取文件名 键
	ConfigNameKey = "CONFIG_NAME"
)

//创建viper.Viper LoadParamsFromEnv将使用env参数
func LoadParamsFromEnv() *viper.Viper {
	v := viper.New()
	prefix := os.Getenv(EnvPrefixKey)
	if prefix == "" {
		prefix = DefaultEnvPrefixKey
		Log.Warn("ENV_PREFIX not exist in env Use default env prefix: %s", DefaultEnvPrefixKey)
	} else {
		Log.Warn("Use EnvPrefixKey: %s", prefix)
	}
	v.SetEnvPrefix(prefix)
	v.AutomaticEnv()
	return v
}

//LoadParamsFromVolume 使用参数创建 viper.Viper
func LoadParamsFromVolume() (*viper.Viper, error) {
	v := viper.New()
	configDir := os.Getenv(ConfigDirKey)
	fileName := os.Getenv(ConfigNameKey)

	//使用默认DIR
	if configDir == "" {
		configDir = DefaultDir
		Log.Warn("ConfigDirKey not exist in env Use default dir %s", DefaultDir)
	} else {
		Log.Info("Use Config_Dir: %s", configDir)
	}

	//使用默认文件名称
	if fileName == "" {
		fileName = DefaultFileName
		Log.Warn("ConfigNameKey not exist in env Use default name %s", DefaultFileName)
	} else {
		Log.Info("Use CONFIG_NAME: %s", fileName)
	}

	v.SetConfigName(fileName)
	v.AddConfigPath(configDir)

	return v, v.ReadInConfig()
}
