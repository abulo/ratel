package redis

import (
	"time"

	"github.com/spf13/viper"
)

// keys of config params from environment variable
var configKey = []string{
	"REDIS_TYPE",
	"REDIS_HOST",
	"REDIS_PORT",
	"REDIS_DB_NAME",
	"REDIS_DB_PASSWORD",
	"REDIS_MAX_CONNECTIONS",
	"REDIS_KEY_PREFIX",
	"REDIS_SKIP_FULL_COVER_CHECK",
	"REDIS_TIMEOUT",
}

// mergeViper will merge tow viper,viper env will cover viper vol
func mergeViper(vol, env *viper.Viper, suffix RWType) {
	for _, value := range configKey {
		value = value + string(suffix)
		if ee := env.Get(value); ee != nil {
			vol.Set(value, env.Get(value))
		}
	}
}

//addrStructure will create ADDR,For example string: "host:port"
func addrStructure(redisPort []string, redisHosts []string) []string {
	hosts := []string{}
	if len(redisPort) != len(redisHosts) {
		port := "6379"
		if len(redisPort) == 0 {
			Log.Warn("REDIS_PORT not exist, Use default port:%s", port)
		} else {
			port = redisPort[0]
			Log.Warn("REDIS_PORT len not equal REDIS_HOST len, Use first port:%s", port)
		}
		for _, host := range redisHosts {
			host := host + ":" + port
			hosts = append(hosts, host)
		}
	} else {
		for index, host := range redisHosts {
			host := host + ":" + redisPort[index]
			hosts = append(hosts, host)
		}
	}
	if len(hosts) == 0 {
		Log.Warn("REDIS_PORT hosts is empty")
	}
	return hosts
}

//customizedOption create options and config the Option
func customizedOption(viper *viper.Viper, rwType RWType) *Options {

	var opt = Options{}
	hosts := addrStructure(viper.GetStringSlice(rwType.FmtSuffix("REDIS_PORT")),
		viper.GetStringSlice(rwType.FmtSuffix("REDIS_HOST")))
	opt.Type = ClientType(viper.GetString(rwType.FmtSuffix("REDIS_TYPE")))
	opt.Hosts = hosts
	opt.ReadOnly = rwType.IsReadOnly()
	opt.Database = viper.GetInt(rwType.FmtSuffix("REDIS_DB_NAME"))
	opt.Password = viper.GetString(rwType.FmtSuffix("REDIS_DB_PASSWORD"))
	opt.KeyPrefix = viper.GetString(rwType.FmtSuffix("REDIS_KEY_PREFIX"))
	// various timeout setting
	opt.DialTimeout = viper.GetDuration(rwType.FmtSuffix("REDIS_TIMEOUT")) * time.Second
	opt.ReadTimeout = viper.GetDuration(rwType.FmtSuffix("REDIS_TIMEOUT")) * time.Second
	opt.WriteTimeout = viper.GetDuration(rwType.FmtSuffix("REDIS_TIMEOUT")) * time.Second
	// REDIS_MAX_CONNECTIONS
	opt.PoolSize = viper.GetInt(rwType.FmtSuffix("REDIS_MAX_CONNECTIONS"))
	opt.PoolTimeout = viper.GetDuration(rwType.FmtSuffix("REDIS_TIMEOUT")) * time.Second
	opt.IdleTimeout = viper.GetDuration(rwType.FmtSuffix("REDIS_TIMEOUT")) * time.Second
	opt.IdleCheckFrequency = viper.GetDuration(rwType.FmtSuffix("REDIS_TIMEOUT")) * time.Second
	opt.TLSConfig = nil
	return &opt
}

// customizedOptionsFromVolume Customized Options by  Volume
func customizedOptionsFromVolume(rwType RWType) (*Options, error) {
	fromVolume, err := LoadParamsFromVolume()
	if err != nil {
		return nil, err
	}
	return customizedOption(fromVolume, rwType), nil
}

// customizedOptionsFromEnv Customized Options by  Env
func customizedOptionsFromEnv(rwType RWType) (*Options, error) {
	fromEnv := LoadParamsFromEnv()
	return customizedOption(fromEnv, rwType), nil
}

// customizedOptionsFromFullVariable Customized Options by  Volume and Env
func customizedOptionsFromFullVariable(rwType RWType) (*Options, error) {

	fromVolume, err := LoadParamsFromVolume()
	if err != nil {
		return nil, err
	}
	fromEnv := LoadParamsFromEnv()
	mergeViper(fromVolume, fromEnv, rwType)
	return customizedOption(fromVolume, rwType), nil
}
