package icron

import (
	"fmt"
	"github.com/xqk/good/pkg/ilog"
	"runtime"
	"time"

	"github.com/xqk/good/pkg/client/etcdv3"
	"github.com/xqk/good/pkg/ecode"
	"go.etcd.io/etcd/client/v3/concurrency"

	"github.com/robfig/cron/v3"
	"github.com/xqk/good/pkg/conf"
	"github.com/xqk/good/pkg/metric"
	"go.uber.org/zap"
)

// StdConfig ...
func StdConfig(name string) Config {
	return RawConfig("good.cron." + name)
}

// RawConfig ...
func RawConfig(key string) Config {
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(key, &config); err != nil {
		ilog.Panic("unmarshal", ilog.String("key", key))
	}

	if config.DistributedTask {
		config.Config = etcdv3.RawConfig(key)
	}

	return config
}

// DefaultConfig ...
func DefaultConfig() Config {
	return Config{
		logger:          ilog.GoodLogger,
		wrappers:        []JobWrapper{},
		WithSeconds:     false,
		ImmediatelyRun:  false,
		ConcurrentDelay: -1, // skip
	}
}

// Config ...
type Config struct {
	WithSeconds     bool
	ConcurrentDelay int
	ImmediatelyRun  bool

	wrappers []JobWrapper
	logger   *ilog.Logger
	parser   cron.Parser

	// Distributed task
	DistributedTask bool

	WaitLockTime time.Duration
	*etcdv3.Config
	client *etcdv3.Client
}

// WithChain ...
func (config *Config) WithChain(wrappers ...JobWrapper) Config {
	if config.wrappers == nil {
		config.wrappers = []JobWrapper{}
	}
	config.wrappers = append(config.wrappers, wrappers...)
	return *config
}

// WithLogger ...
func (config *Config) WithLogger(logger *ilog.Logger) Config {
	config.logger = logger
	return *config
}

// WithParser ...
func (config *Config) WithParser(parser Parser) Config {
	config.parser = parser
	return *config
}

// Build ...
func (config Config) Build() *Cron {
	if config.WithSeconds {
		config.parser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	} else {
		// default parser
		config.parser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	}

	if config.ConcurrentDelay > 0 { // 延迟
		config.wrappers = append(config.wrappers, delayIfStillRunning(config.logger))
	} else if config.ConcurrentDelay < 0 { // 跳过
		config.wrappers = append(config.wrappers, skipIfStillRunning(config.logger))
	} else {
		// 默认不延迟也不跳过
	}

	if config.DistributedTask {
		// 创建 Etcd Lock
		newETCDXcron(&config)
	} else {
		config.Config = &etcdv3.Config{}
	}

	return newCron(&config)
}

func newETCDXcron(config *Config) {
	if config.logger == nil {
		config.logger = ilog.DefaultLogger
	}
	config.logger = config.logger.With(ilog.FieldMod(ecode.ModXcronETCD), ilog.FieldAddrAny(config.Config.Endpoints))
	config.client = config.Config.MustBuild()
	if config.TTL == 0 {
		config.TTL = DefaultTTL
	}

	return
}

type wrappedLogger struct {
	*ilog.Logger
}

// Info logs routine messages about cron's operation.
func (wl *wrappedLogger) Info(msg string, keysAndValues ...interface{}) {
	wl.Infow("cron "+msg, keysAndValues...)
}

// Error logs an error condition.
func (wl *wrappedLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	wl.Errorw("cron "+msg, append(keysAndValues, "err", err)...)
}

type wrappedJob struct {
	NamedJob
	logger *ilog.Logger

	distributedTask bool
	waitLockTime    time.Duration
	leaseTTL        int
	client          *etcdv3.Client
}

const (
	// 任务锁
	WorkerLockDir       = "/xcron/lock/"
	DefaultTTL          = 60   // default set
	DefaultWaitLockTime = 1000 // ms
)

// Run ...
func (wj wrappedJob) Run() {
	if wj.distributedTask {
		mutex, err := wj.client.NewMutex(WorkerLockDir+wj.Name(), concurrency.WithTTL(wj.leaseTTL))
		if err != nil {
			wj.logger.Error("mutex", ilog.String("err", err.Error()))
			return
		}
		if wj.waitLockTime == 0 {
			err = mutex.TryLock(DefaultWaitLockTime * time.Millisecond)
		} else { // 阻塞等待直到waitLockTime timeout
			err = mutex.Lock(wj.waitLockTime)
		}
		if err != nil {
			wj.logger.Info("mutex lock", ilog.String("err", err.Error()))
			return
		}
		defer mutex.Unlock()
	}
	_ = wj.run()
}

func (wj wrappedJob) run() (err error) {
	metric.JobHandleCounter.Inc("cron", wj.Name(), "begin")
	var fields = []ilog.Field{zap.String("name", wj.Name())}
	var beg = time.Now()
	defer func() {
		if rec := recover(); rec != nil {
			switch rec := rec.(type) {
			case error:
				err = rec
			default:
				err = fmt.Errorf("%v", rec)
			}

			stack := make([]byte, 4096)
			length := runtime.Stack(stack, true)
			fields = append(fields, zap.ByteString("stack", stack[:length]))
		}
		if err != nil {
			fields = append(fields, ilog.String("err", err.Error()), ilog.Duration("cost", time.Since(beg)))
			wj.logger.Error("run", fields...)
		} else {
			wj.logger.Info("run", fields...)
		}
		metric.JobHandleHistogram.Observe(time.Since(beg).Seconds(), "cron", wj.Name())
	}()

	return wj.NamedJob.Run()
}
