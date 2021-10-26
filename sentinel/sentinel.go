package sentinel

import (
	"encoding/json"
	"io/ioutil"

	"github.com/abulo/ratel/logger"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	sentinel_config "github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
)

// Config ...
type Config struct {
	AppName       string
	LogPath       string
	FlowRules     []*flow.Rule
	FlowRulesFile string
}

func (config *Config) Build() error {
	if config.FlowRulesFile != "" {
		var rules []*flow.Rule
		content, err := ioutil.ReadFile(config.FlowRulesFile)
		if err != nil {
			logger.Logger.Error("load sentinel flow rules")
		}

		if err := json.Unmarshal(content, &rules); err != nil {
			logger.Logger.Error("load sentinel flow rules")
		}

		config.FlowRules = append(config.FlowRules, rules...)
	}

	configEntity := sentinel_config.NewDefaultConfig()
	configEntity.Sentinel.App.Name = config.AppName
	configEntity.Sentinel.Log.Dir = config.LogPath

	if len(config.FlowRules) > 0 {
		_, _ = flow.LoadRules(config.FlowRules)
	}
	return sentinel.InitWithConfig(configEntity)
}

func Entry(resource string) (*base.SentinelEntry, *base.BlockError) {
	return sentinel.Entry(resource)
}
