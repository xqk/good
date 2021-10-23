package sentinel

import (
	"encoding/json"
	"github.com/xqk/good/pkg/ilog"
	"io/ioutil"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	sentinel_config "github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/xqk/good/pkg"
	"github.com/xqk/good/pkg/conf"
)

const ModuleName = "sentinel"

// StdConfig ...
func StdConfig(name string) *Config {
	return RawConfig("good.sentinel." + name)
}

// RawConfig ...
func RawConfig(key string) *Config {
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(key, config); err != nil {
		ilog.Panic("unmarshal key", ilog.Any("err", err))
	}
	return config
}

// Config ...
type Config struct {
	AppName       string       `json:"appName"`
	LogPath       string       `json:"logPath"`
	FlowRules     []*flow.Rule `json:"rules"`
	FlowRulesFile string       `json:"flowRulesFile"`
}

// DefaultConfig returns default config for sentinel
func DefaultConfig() *Config {
	return &Config{
		AppName:   pkg.Name(),
		LogPath:   "/tmp/log",
		FlowRules: make([]*flow.Rule, 0),
	}
}

// InitSentinelCoreComponent init sentinel core component
// Currently, only flow rules from json file is supported
// todo: support dynamic rule config
// todo: support more rule such as system rule
func (config *Config) Build() error {
	if config.FlowRulesFile != "" {
		var rules []*flow.Rule
		content, err := ioutil.ReadFile(config.FlowRulesFile)
		if err != nil {
			ilog.Error("load sentinel flow rules", ilog.FieldErr(err), ilog.FieldKey(config.FlowRulesFile))
		}

		if err := json.Unmarshal(content, &rules); err != nil {
			ilog.Error("load sentinel flow rules", ilog.FieldErr(err), ilog.FieldKey(config.FlowRulesFile))
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
