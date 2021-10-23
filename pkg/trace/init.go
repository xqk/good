package trace

import (
	"log"

	"github.com/xqk/good/pkg/conf"
	"github.com/xqk/good/pkg/trace/jaeger"
)

func init() {
	// 加载完配置，初始化sentinel
	conf.OnLoaded(func(c *conf.Configuration) {
		log.Print("hook config, init sentinel rules")
		if conf.Get("good.trace.jaeger") != nil {
			var config = jaeger.RawConfig("good.trace.jaeger")
			SetGlobalTracer(config.Build())
		}
	})
}
