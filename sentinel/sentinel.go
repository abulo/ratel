package sentinel

import (
	"encoding/json"
	"io/ioutil"

	"github.com/abulo/ratel/v2/logger"
	"github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
)

// Config ...
type Config struct {
	AppName       string
	LogPath       string
	FlowRules     []*flow.Rule
	FlowRulesFile string
}

func (cfg *Config) Build() error {
	if cfg.FlowRulesFile != "" {
		var rules []*flow.Rule
		content, err := ioutil.ReadFile(cfg.FlowRulesFile)
		if err != nil {
			logger.Logger.Error("load sentinel flow rules")
		}

		if err := json.Unmarshal(content, &rules); err != nil {
			logger.Logger.Error("load sentinel flow rules")
		}

		cfg.FlowRules = append(cfg.FlowRules, rules...)
	}

	configEntity := config.NewDefaultConfig()
	configEntity.Sentinel.App.Name = cfg.AppName
	configEntity.Sentinel.Log.Dir = cfg.LogPath

	if len(cfg.FlowRules) > 0 {
		_, _ = flow.LoadRules(cfg.FlowRules)
	}
	return api.InitWithConfig(configEntity)
}

func Entry(resource string) (*base.SentinelEntry, *base.BlockError) {
	return api.Entry(resource)
}
