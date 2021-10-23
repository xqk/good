package etcdv3

import (
	"github.com/xqk/good/pkg/ilog"
	"time"

	"github.com/xqk/good/pkg/client/etcdv3"
	"github.com/xqk/good/pkg/conf"
	"github.com/xqk/good/pkg/ecode"
	"github.com/xqk/good/pkg/registry"
)

// StdConfig ...
func StdConfig(name string) *Config {
	return RawConfig("good.registry." + name)
}

// RawConfig ...
func RawConfig(key string) *Config {
	var config = DefaultConfig()
	// 解析最外层配置
	if err := conf.UnmarshalKey(key, &config); err != nil {
		ilog.Panic("unmarshal key", ilog.FieldMod("registry.etcd"), ilog.FieldErrKind(ecode.ErrKindUnmarshalConfigErr), ilog.FieldErr(err), ilog.String("key", key), ilog.Any("config", config))
	}
	// 解析嵌套配置
	if err := conf.UnmarshalKey(key, &config.Config); err != nil {
		ilog.Panic("unmarshal key", ilog.FieldMod("registry.etcd"), ilog.FieldErrKind(ecode.ErrKindUnmarshalConfigErr), ilog.FieldErr(err), ilog.String("key", key), ilog.Any("config", config))
	}
	return config
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		Config:      etcdv3.DefaultConfig(),
		ReadTimeout: time.Second * 3,
		Prefix:      "good",
		logger:      ilog.GoodLogger,
		ServiceTTL:  0,
	}
}

// Config ...
type Config struct {
	*etcdv3.Config
	ReadTimeout time.Duration
	ConfigKey   string
	Prefix      string
	ServiceTTL  time.Duration
	logger      *ilog.Logger
}

// Build ...
func (config Config) Build() (registry.Registry, error) {
	if config.ConfigKey != "" {
		config.Config = etcdv3.RawConfig(config.ConfigKey)
	}
	return newETCDRegistry(&config)
}

func (config Config) MustBuild() registry.Registry {
	reg, err := config.Build()
	if err != nil {
		ilog.Panicf("build registry failed: %v", err)
	}
	return reg
}
