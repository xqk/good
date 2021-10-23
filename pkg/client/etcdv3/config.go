package etcdv3

import (
	"github.com/xqk/good/pkg/ilog"
	"github.com/xqk/good/pkg/util/itime"
	"time"

	"github.com/xqk/good/pkg/conf"
	"github.com/xqk/good/pkg/constant"
	"github.com/xqk/good/pkg/ecode"
	"github.com/xqk/good/pkg/flag"
)

var ConfigPrefix = constant.ConfigPrefix + ".etcdv3"

// Config ...
type (
	Config struct {
		Endpoints []string `json:"endpoints"`
		CertFile  string   `json:"certFile"`
		KeyFile   string   `json:"keyFile"`
		CaCert    string   `json:"caCert"`
		BasicAuth bool     `json:"basicAuth"`
		UserName  string   `json:"userName"`
		Password  string   `json:"-"`
		// 连接超时时间
		ConnectTimeout time.Duration `json:"connectTimeout"`
		Secure         bool          `json:"secure"`
		// 自动同步member list的间隔
		AutoSyncInterval time.Duration `json:"autoAsyncInterval"`
		TTL              int           // 单位：s
		logger           *ilog.Logger
	}
)

func (config *Config) BindFlags(fs *flag.FlagSet) {
	fs.BoolVar(&config.Secure, "insecure-etcd", true, "--insecure-etcd=true")
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		BasicAuth:      false,
		ConnectTimeout: itime.Duration("5s"),
		Secure:         false,
		logger:         ilog.GoodLogger.With(ilog.FieldMod("client.etcd")),
	}
}

// StdConfig ...
func StdConfig(name string) *Config {
	return RawConfig(ConfigPrefix + "." + name)
}

// RawConfig ...
func RawConfig(key string) *Config {
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(key, config); err != nil {
		config.logger.Panic("client etcd parse config panic", ilog.FieldErrKind(ecode.ErrKindUnmarshalConfigErr), ilog.FieldErr(err), ilog.FieldKey(key), ilog.FieldValueAny(config))
	}
	return config
}

// WithLogger ...
func (config *Config) WithLogger(logger *ilog.Logger) *Config {
	config.logger = logger
	return config
}

// Build ...
func (config *Config) Build() (*Client, error) {
	return newClient(config)
}

func (config *Config) MustBuild() *Client {
	client, err := config.Build()
	if err != nil {
		ilog.Panicf("build etcd client failed: %v", err)
	}
	return client
}
