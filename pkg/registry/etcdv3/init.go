package etcdv3

import (
	"github.com/xqk/good/pkg/registry"
)

func init() {
	registry.RegisterBuilder("etcdv3", func(confKey string) registry.Registry {
		return RawConfig(confKey).MustBuild()
	})
}
