package autoproc

import (
	"github.com/xqk/good/pkg/conf"
	"github.com/xqk/good/pkg/ecode"
	"github.com/xqk/good/pkg/ilog"
	"go.uber.org/automaxprocs/maxprocs"
	"runtime"
)

func init() {
	// 初始化注册中心
	if _, err := maxprocs.Set(); err != nil {
		ilog.Panic("auto max procs", ilog.FieldMod(ecode.ModProc), ilog.FieldErrKind(ecode.ErrKindAny), ilog.FieldErr(err))
	}
	conf.OnLoaded(func(c *conf.Configuration) {
		if maxProcs := conf.GetInt("maxProc"); maxProcs != 0 {
			runtime.GOMAXPROCS(maxProcs)
		}
		ilog.Info("auto max procs", ilog.FieldMod(ecode.ModProc), ilog.Int64("procs", int64(runtime.GOMAXPROCS(-1))))
	})
}
