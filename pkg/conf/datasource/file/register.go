package file

import (
	"github.com/xqk/good/pkg/conf"
	"github.com/xqk/good/pkg/flag"
	"github.com/xqk/good/pkg/ilog"
)

// DataSourceFile defines file scheme
const DataSourceFile = "file"

func init() {
	conf.Register(DataSourceFile, func() conf.DataSource {
		var (
			watchConfig = flag.Bool("watch")
			configAddr  = flag.String("config")
		)
		if configAddr == "" {
			ilog.Panic("new file dataSource, configAddr is empty")
			return nil
		}
		return NewDataSource(configAddr, watchConfig)
	})
}
