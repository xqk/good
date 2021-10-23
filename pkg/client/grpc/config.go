package grpc

import (
	"github.com/xqk/good/pkg/conf"
	"github.com/xqk/good/pkg/ecode"
	"github.com/xqk/good/pkg/ilog"
	"github.com/xqk/good/pkg/util/itime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/keepalive"
	"time"
)

// Config ...
type Config struct {
	Name         string // config's name
	BalancerName string
	Address      string
	Block        bool
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	Direct       bool
	OnDialError  string // panic | error
	KeepAlive    *keepalive.ClientParameters
	logger       *ilog.Logger
	dialOptions  []grpc.DialOption

	SlowThreshold time.Duration

	Debug                     bool
	DisableTraceInterceptor   bool
	DisableAidInterceptor     bool
	DisableTimeoutInterceptor bool
	DisableMetricInterceptor  bool
	DisableAccessInterceptor  bool
	AccessInterceptorLevel    string
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		dialOptions: []grpc.DialOption{
			grpc.WithInsecure(),
		},
		logger:                 ilog.GoodLogger.With(ilog.FieldMod(ecode.ModClientGrpc)),
		BalancerName:           roundrobin.Name, // round robin by default
		DialTimeout:            time.Second * 3,
		ReadTimeout:            itime.Duration("1s"),
		SlowThreshold:          itime.Duration("600ms"),
		OnDialError:            "panic",
		AccessInterceptorLevel: "info",
		Block:                  true,
	}
}

// StdConfig ...
func StdConfig(name string) *Config {
	return RawConfig("good.client." + name)
}

// RawConfig ...
func RawConfig(key string) *Config {
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(key, &config); err != nil {
		config.logger.Panic("client grpc parse config panic", ilog.FieldErrKind(ecode.ErrKindUnmarshalConfigErr), ilog.FieldErr(err), ilog.FieldKey(key), ilog.FieldValueAny(config))
	}
	return config
}

// WithLogger ...
func (config *Config) WithLogger(logger *ilog.Logger) *Config {
	config.logger = logger
	return config
}

// WithDialOption ...
func (config *Config) WithDialOption(opts ...grpc.DialOption) *Config {
	if config.dialOptions == nil {
		config.dialOptions = make([]grpc.DialOption, 0)
	}
	config.dialOptions = append(config.dialOptions, opts...)
	return config
}

// Build ...
func (config *Config) Build() *grpc.ClientConn {
	if config.Debug {
		config.dialOptions = append(config.dialOptions,
			grpc.WithChainUnaryInterceptor(debugUnaryClientInterceptor(config.Address)),
		)
	}

	if !config.DisableAidInterceptor {
		config.dialOptions = append(config.dialOptions,
			grpc.WithChainUnaryInterceptor(aidUnaryClientInterceptor()),
		)
	}

	if !config.DisableTimeoutInterceptor {
		config.dialOptions = append(config.dialOptions,
			grpc.WithChainUnaryInterceptor(timeoutUnaryClientInterceptor(config.logger, config.ReadTimeout, config.SlowThreshold)),
		)
	}

	if !config.DisableTraceInterceptor {
		config.dialOptions = append(config.dialOptions,
			grpc.WithChainUnaryInterceptor(traceUnaryClientInterceptor()),
		)
	}

	if !config.DisableAccessInterceptor {
		config.dialOptions = append(config.dialOptions,
			grpc.WithChainUnaryInterceptor(loggerUnaryClientInterceptor(config.logger, config.Name, config.AccessInterceptorLevel)),
		)
	}

	if !config.DisableMetricInterceptor {
		config.dialOptions = append(config.dialOptions,
			grpc.WithChainUnaryInterceptor(metricUnaryClientInterceptor(config.Name)),
		)
	}

	return newGRPCClient(config)
}
