package http

import (
	"github.com/xqk/good/pkg/conf"
	"github.com/xqk/good/pkg/flag"
	"github.com/xqk/good/pkg/ilog"
)

// Defines http/https scheme
const (
	DataSourceHttp  = "http"
	DataSourceHttps = "https"
)

func init() {
	dataSourceCreator := func() conf.DataSource {
		var (
			watchConfig = flag.Bool("watch")
			configAddr  = flag.String("config")
		)
		if configAddr == "" {
			ilog.Panic("new http dataSource, configAddr is empty")
			return nil
		}
		return NewDataSource(configAddr, watchConfig)
	}
	conf.Register(DataSourceHttp, dataSourceCreator)
	conf.Register(DataSourceHttps, dataSourceCreator)
}
