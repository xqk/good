package gorm

import (
	"github.com/xqk/good/pkg/conf"
	"github.com/xqk/good/pkg/ecode"
	"github.com/xqk/good/pkg/ilog"
	"github.com/xqk/good/pkg/metric"
	"github.com/xqk/good/pkg/util/itime"
	"time"
)

// StdConfig 标准配置，规范配置文件头
func StdConfig(name string) *Config {
	return RawConfig("good.mysql." + name)
}

// RawConfig 传入mapstructure格式的配置
// example: RawConfig("good.mysql.stt_config")
func RawConfig(key string) *Config {
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(key, config, conf.TagName("toml")); err != nil {
		ilog.Panic("unmarshal key", ilog.FieldMod("gorm"), ilog.FieldErr(err), ilog.FieldKey(key))
	}
	config.Name = key
	return config
}

// config options
type Config struct {
	Name string
	// DSN地址: mysql://root:secret@tcp(127.0.0.1:3307)/mysql?timeout=20s&readTimeout=20s
	DSN string `json:"dsn" toml:"dsn"`
	// Debug开关
	Debug bool `json:"debug" toml:"debug"`
	// 最大空闲连接数
	MaxIdleConns int `json:"maxIdleConns" toml:"maxIdleConns"`
	// 最大活动连接数
	MaxOpenConns int `json:"maxOpenConns" toml:"maxOpenConns"`
	// 连接的最大存活时间
	ConnMaxLifetime time.Duration `json:"connMaxLifetime" toml:"connMaxLifetime"`
	// 创建连接的错误级别，=panic时，如果创建失败，立即panic
	OnDialError string `json:"level" toml:"level"`
	// 慢日志阈值
	SlowThreshold time.Duration `json:"slowThreshold" toml:"slowThreshold"`
	// 拨超时时间
	DialTimeout time.Duration `json:"dialTimeout" toml:"dialTimeout"`
	// 关闭指标采集
	DisableMetric bool `json:"disableMetric" toml:"disableMetric"`
	// 关闭链路追踪
	DisableTrace bool `json:"disableTrace" toml:"disableTrace"`

	// 记录错误sql时,是否打印包含参数的完整sql语句
	// select * from aid = ?;
	// select * from aid = 288016;
	DetailSQL bool `json:"detailSql" toml:"detailSql"`

	raw          interface{}
	logger       *ilog.Logger
	interceptors []Interceptor
	dsnCfg       *DSN
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		DSN:             "",
		Debug:           false,
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: itime.Duration("300s"),
		OnDialError:     "panic",
		SlowThreshold:   itime.Duration("500ms"),
		DialTimeout:     itime.Duration("1s"),
		DisableMetric:   false,
		DisableTrace:    false,
		raw:             nil,
		logger:          ilog.GoodLogger,
	}
}

// WithLogger ...
func (config *Config) WithLogger(log *ilog.Logger) *Config {
	config.logger = log
	return config
}

// WithInterceptor ...
func (config *Config) WithInterceptor(intes ...Interceptor) *Config {
	if config.interceptors == nil {
		config.interceptors = make([]Interceptor, 0)
	}
	config.interceptors = append(config.interceptors, intes...)
	return config
}

// Build ...
func (config *Config) Build() *DB {
	var err error
	config.dsnCfg, err = ParseDSN(config.DSN)
	if err == nil {
		config.logger.Info(ecode.MsgClientMysqlOpenStart, ilog.FieldMod("gorm"), ilog.FieldAddr(config.dsnCfg.Addr), ilog.FieldName(config.dsnCfg.DBName))
	} else {
		config.logger.Panic(ecode.MsgClientMysqlOpenStart, ilog.FieldMod("gorm"), ilog.FieldErr(err))
	}

	if config.Debug {
		config = config.WithInterceptor(debugInterceptor)
	}
	if !config.DisableTrace {
		config = config.WithInterceptor(traceInterceptor)
	}

	if !config.DisableMetric {
		config = config.WithInterceptor(metricInterceptor)
	}

	db, err := Open("mysql", config)
	if err != nil {
		if config.OnDialError == "panic" {
			config.logger.Panic("open mysql", ilog.FieldMod("gorm"), ilog.FieldErrKind(ecode.ErrKindRequestErr), ilog.FieldErr(err), ilog.FieldAddr(config.dsnCfg.Addr), ilog.FieldValueAny(config))
		} else {
			metric.LibHandleCounter.Inc(metric.TypeGorm, config.Name+".ping", config.dsnCfg.Addr, "open err")
			config.logger.Error("open mysql", ilog.FieldMod("gorm"), ilog.FieldErrKind(ecode.ErrKindRequestErr), ilog.FieldErr(err), ilog.FieldAddr(config.dsnCfg.Addr), ilog.FieldValueAny(config))
			return db
		}
	}

	if err := db.DB().Ping(); err != nil {
		config.logger.Panic("ping mysql", ilog.FieldMod("gorm"), ilog.FieldErrKind(ecode.ErrKindRequestErr), ilog.FieldErr(err), ilog.FieldValueAny(config))
	}

	// store db
	instances.Store(config.Name, db)
	return db
}
