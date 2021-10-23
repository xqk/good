package governor

import (
	"fmt"
	"github.com/xqk/good/pkg/ilog"

	"github.com/xqk/good/pkg/conf"
	"github.com/xqk/good/pkg/util/inet"
)

//ModName ..
const ModName = "govern"

// Config ...
type Config struct {
	Host    string
	Port    int
	Network string `json:"network" toml:"network"`
	logger  *ilog.Logger
	Enable  bool

	// ServiceAddress service address in registry info, default to 'Host:Port'
	ServiceAddress string
}

// StdConfig represents Standard gRPC Server config
// which will parse config by conf package,
// panic if no config key found in conf
func StdConfig(name string) *Config {
	return RawConfig("good.server." + name)
}

// RawConfig ...
func RawConfig(key string) *Config {
	var config = DefaultConfig()
	if conf.Get(key) == nil {
		return config
	}
	if err := conf.UnmarshalKey(key, &config); err != nil {
		config.logger.Panic("govern server parse config panic",
			ilog.FieldErr(err), ilog.FieldKey(key),
			ilog.FieldValueAny(config),
		)
	}
	return config
}

// DefaultConfig represents default config
// User should construct config base on DefaultConfig
func DefaultConfig() *Config {
	host, port, err := inet.GetLocalMainIP()
	if err != nil {
		host = "localhost"
	}

	return &Config{
		Enable:  true,
		Host:    host,
		Network: "tcp4",
		Port:    port,
		logger:  ilog.GoodLogger.With(ilog.FieldMod(ModName)),
	}
}

// Build ...
func (config *Config) Build() *Server {
	return newServer(config)
}

// Address ...
func (config Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}
